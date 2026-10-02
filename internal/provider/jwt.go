package provider

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	// jwtExpirySkew makes the provider renew a JWT shortly before it expires so
	// that a long apply does not start with a token that dies mid-way.
	jwtExpirySkew = 10 * time.Minute

	// jwtFallbackLifetime is assumed for a JWT without an exp claim.
	jwtFallbackLifetime = time.Hour

	maxSignInResponse = 1 << 16
)

type apiCredentials struct {
	email    string
	apiToken string
}

func (c apiCredentials) complete() bool { return c.email != "" && c.apiToken != "" }

// tokenSource hands out the Elestio session JWT. The backend limits sign-ins
// (15 per user and per IP per hour) but issues JWTs that stay valid for days,
// so the token is reused until it nears expiry instead of signing in on every
// Terraform command. It is safe for concurrent use.
type tokenSource struct {
	signInURL string
	creds     apiCredentials
	client    *http.Client // used for sign-in only; carries no auth transport
	cache     *jwtCache    // optional
	now       func() time.Time

	mu    sync.Mutex
	token string
	exp   time.Time
}

func newTokenSource(baseURL string, creds apiCredentials, client *http.Client, cache *jwtCache) *tokenSource {
	return &tokenSource{
		signInURL: strings.TrimRight(baseURL, "/") + "/api/auth/checkAPIToken",
		creds:     creds,
		client:    client,
		cache:     cache,
		now:       time.Now,
	}
}

func (s *tokenSource) usable(token string, exp time.Time) bool {
	return token != "" && s.now().Add(jwtExpirySkew).Before(exp)
}

// seed installs a JWT supplied by the user (ELESTIO_JWT). An expired or
// malformed one is ignored.
func (s *tokenSource) seed(token string) {
	exp, ok := jwtExpiry(token, s.now())
	if !ok || !s.usable(token, exp) {
		return
	}
	s.mu.Lock()
	s.token, s.exp = token, exp
	s.mu.Unlock()
}

// Token returns a JWT that is not about to expire, using in order: the
// in-memory token, the on-disk cache, and a fresh sign-in.
func (s *tokenSource) Token(ctx context.Context) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.usable(s.token, s.exp) {
		return s.token, nil
	}
	if s.cache != nil {
		if tok, exp, ok := s.cache.load(s.now()); ok && s.usable(tok, exp) {
			s.token, s.exp = tok, exp
			return tok, nil
		}
	}
	return s.signInLocked(ctx)
}

// Refresh replaces a token the API rejected as invalid. Callers pass the token
// they used; if another goroutine already replaced it, that newer token is
// returned without a second sign-in.
func (s *tokenSource) Refresh(ctx context.Context, rejected string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.token != "" && s.token != rejected && s.usable(s.token, s.exp) {
		return s.token, nil
	}
	s.token, s.exp = "", time.Time{}
	if s.cache != nil {
		s.cache.remove()
	}
	return s.signInLocked(ctx)
}

func (s *tokenSource) signInLocked(ctx context.Context) (string, error) {
	if !s.creds.complete() {
		return "", errors.New("the Elestio JWT is missing, expired or was rejected, and no email and api_token are configured to sign in again")
	}

	payload, err := json.Marshal(map[string]string{"email": s.creds.email, "token": s.creds.apiToken})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.signInURL, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to sign in: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxSignInResponse))
	if err != nil {
		return "", fmt.Errorf("failed to sign in: reading response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("failed to sign in: request failed with status code %d: %s", resp.StatusCode, body)
	}

	var out struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		JWT     string `json:"jwt"`
	}
	if err := json.Unmarshal(body, &out); err != nil {
		return "", errors.New("failed to sign in: the API returned a response that is not valid JSON")
	}
	if out.Status == "KO" {
		return "", fmt.Errorf("failed to sign in: %s", out.Message)
	}
	if out.JWT == "" {
		return "", errors.New("failed to sign in: the API response contains no JWT")
	}

	exp, ok := jwtExpiry(out.JWT, s.now())
	if !ok {
		exp = s.now().Add(jwtFallbackLifetime)
	}
	s.token, s.exp = out.JWT, exp
	if s.cache != nil {
		s.cache.store(out.JWT, exp)
	}
	return out.JWT, nil
}

// jwtExpiry reads the exp claim of a JWT without verifying the signature. The
// result only decides when to renew locally; the API remains the authority on
// whether a token is valid.
func jwtExpiry(token string, now time.Time) (time.Time, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return time.Time{}, false
	}
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(parts[1], "="))
	if err != nil {
		return time.Time{}, false
	}
	var claims struct {
		Exp float64 `json:"exp"`
	}
	if json.Unmarshal(raw, &claims) != nil || claims.Exp <= 0 {
		return time.Time{}, false
	}
	return time.Unix(int64(claims.Exp), 0), true
}

// ---------------------------------------------------------------- disk cache

// jwtCache stores the JWT between Terraform runs so that plan, apply and
// refresh do not each spend one of the hourly sign-ins. It is OFF by default,
// because it leaves a credential that stays valid for days on disk. Set
// ELESTIO_JWT_CACHE=on to enable it. The file is private to the user (0600 in a
// 0700 directory), keyed by a hash of the credentials, and removed when the API
// rejects the token.
type jwtCache struct{ path string }

func newJWTCache(creds apiCredentials) *jwtCache {
	switch strings.ToLower(os.Getenv("ELESTIO_JWT_CACHE")) {
	case "on", "1", "true", "yes", "enabled":
	default:
		return nil
	}
	if !creds.complete() {
		return nil
	}
	dir, err := os.UserCacheDir()
	if err != nil || dir == "" {
		return nil
	}
	sum := sha256.Sum256([]byte(creds.email + "\x00" + creds.apiToken))
	return &jwtCache{path: filepath.Join(dir, "terraform-provider-elestio", hex.EncodeToString(sum[:16])+".json")}
}

type cachedJWT struct {
	JWT string `json:"jwt"`
	Exp int64  `json:"exp"`
}

func (c *jwtCache) load(now time.Time) (string, time.Time, bool) {
	info, err := os.Lstat(c.path)
	if err != nil || !info.Mode().IsRegular() {
		return "", time.Time{}, false
	}
	// A cache readable by other users is not trusted (file permissions are not
	// meaningful on Windows).
	if runtime.GOOS != "windows" && info.Mode().Perm()&0o077 != 0 {
		_ = os.Remove(c.path)
		return "", time.Time{}, false
	}
	raw, err := os.ReadFile(c.path)
	if err != nil {
		return "", time.Time{}, false
	}
	var v cachedJWT
	if json.Unmarshal(raw, &v) != nil || v.JWT == "" {
		return "", time.Time{}, false
	}
	exp := time.Unix(v.Exp, 0)
	if !now.Before(exp) {
		return "", time.Time{}, false
	}
	return v.JWT, exp, true
}

// store writes the cache atomically. Failures are ignored: the cache is an
// optimization and must never break a run.
func (c *jwtCache) store(token string, exp time.Time) {
	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return
	}
	// #nosec G117 -- persisting the session JWT is the purpose of this cache. It is
	// opt-in (ELESTIO_JWT_CACHE=on), written 0600 in a 0700 directory, keyed by a
	// credential hash, and never logged (see docs/index.md).
	raw, err := json.Marshal(cachedJWT{JWT: token, Exp: exp.Unix()})
	if err != nil {
		return
	}
	tmp, err := os.CreateTemp(dir, ".jwt-*")
	if err != nil {
		return
	}
	name := tmp.Name()
	defer func() { _ = os.Remove(name) }()
	if err := tmp.Chmod(0o600); err != nil {
		_ = tmp.Close()
		return
	}
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	_ = os.Rename(name, c.path)
}

func (c *jwtCache) remove() { _ = os.Remove(c.path) }
