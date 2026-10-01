package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func makeJWT(t *testing.T, exp time.Time, id string) string {
	t.Helper()
	enc := func(v any) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	return enc(map[string]string{"alg": "HS256"}) + "." + enc(map[string]any{"exp": exp.Unix(), "jti": id}) + ".c2ln"
}

// fakeSignIn is a stand-in for /api/auth/checkAPIToken that issues a fresh JWT
// per call and counts how many it issued.
type fakeSignIn struct {
	t      *testing.T
	calls  atomic.Int32
	life   time.Duration
	status int
	body   string // overrides the response when set
}

func (f *fakeSignIn) token(n int32) string {
	life := f.life
	if life == 0 {
		life = 7 * 24 * time.Hour
	}
	return makeJWT(f.t, time.Now().Add(life), fmt.Sprintf("tok-%d", n))
}

func (f *fakeSignIn) handler(w http.ResponseWriter, r *http.Request) {
	n := f.calls.Add(1)
	if f.status != 0 {
		w.WriteHeader(f.status)
	}
	if f.body != "" {
		_, _ = w.Write([]byte(f.body))
		return
	}
	_, _ = fmt.Fprintf(w, `{"status":"OK","jwt":%q}`, f.token(n))
}

func newSource(t *testing.T, url string, creds apiCredentials, cache *jwtCache) *tokenSource {
	t.Helper()
	return newTokenSource(url, creds, http.DefaultClient, cache)
}

var testCreds = apiCredentials{email: "user@example.com", apiToken: "eAt_secret_api_token"}

func TestTokenSource_ReusesTokenInsteadOfSigningInEveryCall(t *testing.T) {
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()
	src := newSource(t, srv.URL, testCreds, nil)

	first, err := src.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 20; i++ {
		got, err := src.Token(context.Background())
		if err != nil || got != first {
			t.Fatalf("call %d: token changed or failed: %v", i, err)
		}
	}
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("signed in %d times, want exactly 1", n)
	}
}

func TestTokenSource_RenewsShortlyBeforeExpiry(t *testing.T) {
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()
	src := newSource(t, srv.URL, testCreds, nil)

	now := time.Now()
	src.now = func() time.Time { return now }
	if _, err := src.Token(context.Background()); err != nil {
		t.Fatal(err)
	}

	now = now.Add(7*24*time.Hour - jwtExpirySkew - time.Minute) // still outside the skew
	if _, err := src.Token(context.Background()); err != nil || f.calls.Load() != 1 {
		t.Fatalf("renewed too early (calls=%d, err=%v)", f.calls.Load(), err)
	}
	now = now.Add(2 * time.Minute) // now inside the skew window
	if _, err := src.Token(context.Background()); err != nil || f.calls.Load() != 2 {
		t.Fatalf("did not renew near expiry (calls=%d, err=%v)", f.calls.Load(), err)
	}
}

func TestTokenSource_ConcurrentCallsSignInOnce(t *testing.T) {
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()
	src := newSource(t, srv.URL, testCreds, nil)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := src.Token(context.Background()); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("50 concurrent callers caused %d sign-ins, want 1", n)
	}
}

func TestTokenSource_SignInFailuresAreClearAndDoNotLeak(t *testing.T) {
	cases := map[string]*fakeSignIn{
		"rate limited":    {status: 429, body: `{"status":"KO","code":"TooManyRequests","message":"Access temporarily restricted."}`},
		"bad credentials": {status: 401, body: `{"status":"KO","code":"InvalidCredentials"}`},
		"KO with 200":     {body: `{"status":"KO","message":"nope"}`},
		"no jwt":          {body: `{"status":"OK"}`},
		"not json":        {body: `<html>gateway</html>`},
	}
	for name, f := range cases {
		t.Run(name, func(t *testing.T) {
			f.t = t
			srv := httptest.NewServer(http.HandlerFunc(f.handler))
			defer srv.Close()
			_, err := newSource(t, srv.URL, testCreds, nil).Token(context.Background())
			if err == nil {
				t.Fatal("expected an error")
			}
			if strings.Contains(err.Error(), testCreds.apiToken) {
				t.Fatalf("error leaks the api token: %v", err)
			}
		})
	}
}

func TestTokenSource_RefreshDoesNotSignInTwiceForOneRejection(t *testing.T) {
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()
	src := newSource(t, srv.URL, testCreds, nil)

	old, _ := src.Token(context.Background())
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := src.Refresh(context.Background(), old); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	if n := f.calls.Load(); n != 2 {
		t.Fatalf("sign-ins = %d, want 2 (initial + one refresh)", n)
	}
}

func TestTokenSource_UserSuppliedJWT(t *testing.T) {
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()

	valid := makeJWT(t, time.Now().Add(48*time.Hour), "user")
	src := newSource(t, srv.URL, apiCredentials{}, nil) // JWT only, no credentials
	src.seed(valid)
	got, err := src.Token(context.Background())
	if err != nil || got != valid || f.calls.Load() != 0 {
		t.Fatalf("a valid supplied JWT must be used without signing in (calls=%d err=%v)", f.calls.Load(), err)
	}

	expired := newSource(t, srv.URL, apiCredentials{}, nil)
	expired.seed(makeJWT(t, time.Now().Add(-time.Hour), "old"))
	if _, err := expired.Token(context.Background()); err == nil || !strings.Contains(err.Error(), "no email and api_token") {
		t.Fatalf("an expired JWT with no credentials must fail clearly, got %v", err)
	}

	withCreds := newSource(t, srv.URL, testCreds, nil)
	withCreds.seed("garbage")
	if _, err := withCreds.Token(context.Background()); err != nil || f.calls.Load() != 1 {
		t.Fatalf("a malformed JWT must fall back to signing in (calls=%d err=%v)", f.calls.Load(), err)
	}
}

func TestJWTExpiry(t *testing.T) {
	exp := time.Now().Add(time.Hour).Truncate(time.Second)
	got, ok := jwtExpiry(makeJWT(t, exp, "x"), time.Now())
	if !ok || !got.Equal(exp) {
		t.Fatalf("got %v %v", got, ok)
	}
	for _, bad := range []string{"", "a.b", "a.b.c", "a.!!!.c", makeJWT(t, time.Time{}, "x")} {
		if _, ok := jwtExpiry(bad, time.Now()); ok {
			t.Errorf("%q must be rejected", bad)
		}
	}
}

// ------------------------------------------------------------------ disk cache

func cacheEnv(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("cache location differs on Windows")
	}
	dir := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", dir)
	t.Setenv("HOME", dir)
	t.Setenv("ELESTIO_JWT_CACHE", "on")
	return dir
}

func TestJWTCache_SharedBetweenRuns(t *testing.T) {
	cacheEnv(t)
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()

	// Two separate token sources stand in for two terraform processes.
	for run := 1; run <= 3; run++ {
		src := newSource(t, srv.URL, testCreds, newJWTCache(testCreds))
		if _, err := src.Token(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	if n := f.calls.Load(); n != 1 {
		t.Fatalf("3 runs signed in %d times, want 1", n)
	}
}

func TestJWTCache_FilePermissionsAndKeying(t *testing.T) {
	cacheEnv(t)
	c := newJWTCache(testCreds)
	c.store(makeJWT(t, time.Now().Add(time.Hour), "a"), time.Now().Add(time.Hour))

	fi, err := os.Stat(c.path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("cache file mode %v, want 0600", fi.Mode().Perm())
	}
	di, _ := os.Stat(filepath.Dir(c.path))
	if di.Mode().Perm() != 0o700 {
		t.Fatalf("cache dir mode %v, want 0700", di.Mode().Perm())
	}
	if strings.Contains(c.path, "user@example.com") || strings.Contains(c.path, testCreds.apiToken) {
		t.Fatal("cache file name must not reveal credentials")
	}
	other := newJWTCache(apiCredentials{email: "user@example.com", apiToken: "different"})
	if other.path == c.path {
		t.Fatal("different credentials must use a different cache file")
	}
}

func TestJWTCache_RejectsUnsafeOrStaleFiles(t *testing.T) {
	cacheEnv(t)
	c := newJWTCache(testCreds)
	good := makeJWT(t, time.Now().Add(time.Hour), "g")

	c.store(good, time.Now().Add(time.Hour))
	if _, _, ok := c.load(time.Now()); !ok {
		t.Fatal("fresh cache should load")
	}

	if err := os.Chmod(c.path, 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, ok := c.load(time.Now()); ok {
		t.Fatal("a world-readable cache must not be trusted")
	}
	if _, err := os.Stat(c.path); err == nil {
		t.Fatal("an untrusted cache file should be deleted")
	}

	c.store(good, time.Now().Add(-time.Minute))
	if _, _, ok := c.load(time.Now()); ok {
		t.Fatal("an expired cache must be ignored")
	}

	_ = os.WriteFile(c.path, []byte("not json"), 0o600)
	if _, _, ok := c.load(time.Now()); ok {
		t.Fatal("a corrupt cache must be ignored")
	}
}

func TestJWTCache_IsOffByDefaultAndOptIn(t *testing.T) {
	cacheEnv(t)
	for _, v := range []string{"", "off", "OFF", "0", "false", "no", "garbage"} {
		t.Setenv("ELESTIO_JWT_CACHE", v)
		if newJWTCache(testCreds) != nil {
			t.Errorf("ELESTIO_JWT_CACHE=%q must NOT enable the cache", v)
		}
	}
	for _, v := range []string{"on", "ON", "1", "true", "yes"} {
		t.Setenv("ELESTIO_JWT_CACHE", v)
		if newJWTCache(testCreds) == nil {
			t.Errorf("ELESTIO_JWT_CACHE=%q must enable the cache", v)
		}
	}
	t.Setenv("ELESTIO_JWT_CACHE", "on")
	if newJWTCache(apiCredentials{}) != nil {
		t.Error("no cache without credentials")
	}
}

// With the default settings nothing is written to disk, and a second run signs
// in again (in-memory reuse only lasts within one process).
func TestDefaultSettingsWriteNothingToDisk(t *testing.T) {
	dir := cacheEnv(t)
	t.Setenv("ELESTIO_JWT_CACHE", "")
	f := &fakeSignIn{t: t}
	srv := httptest.NewServer(http.HandlerFunc(f.handler))
	defer srv.Close()

	for i := 0; i < 2; i++ {
		client, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
		if err != nil || client == nil {
			t.Fatal(err)
		}
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 0 {
		t.Fatalf("default settings must not write a cache, found %v", entries)
	}
	if n := f.calls.Load(); n != 2 {
		t.Fatalf("two separate runs signed in %d times, want 2 (no cache)", n)
	}
}

// --------------------------------------------- transport, end to end through the client

// fakeAPIServer implements sign-in plus one authenticated endpoint. It records
// where each request carried the JWT.
type fakeAPIServer struct {
	t         *testing.T
	signIn    fakeSignIn
	mu        sync.Mutex
	seen      []seenReq
	rejectN   int32 // reject this many authenticated calls with the given body
	rejectRaw string
	apiCalls  atomic.Int32
}

type seenReq struct{ header, query, body, payload string }

func (a *fakeAPIServer) handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/auth/checkAPIToken":
		a.signIn.handler(w, r)
	case "/api/projects/getList":
		b, _ := io.ReadAll(r.Body)
		var body map[string]any
		_ = json.Unmarshal(b, &body)
		bodyJWT, _ := body["jwt"].(string)
		a.mu.Lock()
		a.seen = append(a.seen, seenReq{
			header:  strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer "),
			query:   r.URL.Query().Get("jwt"),
			body:    bodyJWT,
			payload: string(b),
		})
		a.mu.Unlock()
		n := a.apiCalls.Add(1)
		if n <= a.rejectN {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(a.rejectRaw))
			return
		}
		_, _ = w.Write([]byte(`{"status":"OK","data":{"projects":[{"projectID":"12345","projectName":"terraform"}]}}`))
	default:
		http.NotFound(w, r)
	}
}

func newFakeAPI(t *testing.T) (*fakeAPIServer, *httptest.Server) {
	a := &fakeAPIServer{t: t, signIn: fakeSignIn{t: t}}
	srv := httptest.NewServer(http.HandlerFunc(a.handler))
	t.Cleanup(srv.Close)
	return a, srv
}

func TestAuthTransport_InjectsJWTEverywhereAndReusesIt(t *testing.T) {
	cacheEnv(t)
	t.Setenv("ELESTIO_JWT_CACHE", "off")
	api, srv := newFakeAPI(t)

	client, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := client.Project.GetList(); err != nil {
			t.Fatalf("call %d: %v", i, err)
		}
	}

	if n := api.signIn.calls.Load(); n != 1 {
		t.Fatalf("5 API calls caused %d sign-ins, want 1", n)
	}
	want := api.signIn.token(1)
	for i, s := range api.seen {
		if s.header != want || s.query != want || s.body != want {
			t.Fatalf("request %d did not carry the JWT in all three places: header=%t query=%t body=%t",
				i, s.header == want, s.query == want, s.body == want)
		}
	}
}

func TestAuthTransport_ExpiredTokenIsRefreshedOnceAndRequestReplayed(t *testing.T) {
	cacheEnv(t)
	t.Setenv("ELESTIO_JWT_CACHE", "off")
	api, srv := newFakeAPI(t)
	api.rejectN = 1
	api.rejectRaw = `{"status":"KO","code":"Unauthorized","message":"Authentication failed."}`

	client, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Project.GetList(); err != nil {
		t.Fatalf("the request must succeed after a transparent refresh: %v", err)
	}
	if n := api.signIn.calls.Load(); n != 2 {
		t.Fatalf("sign-ins = %d, want 2 (initial + refresh)", n)
	}
	if len(api.seen) != 2 || api.seen[0].header == api.seen[1].header {
		t.Fatalf("the replay must use the new token: %+v", api.seen)
	}
	if api.seen[1].body != api.signIn.token(2) {
		t.Fatal("replayed body does not carry the refreshed JWT")
	}
}

func TestAuthTransport_OnlyUnauthorizedCodeTriggersRefresh(t *testing.T) {
	for _, code := range []string{"InvalidServer", "service_deleted", "InvalidPermission", "MissingPermission"} {
		t.Run(code, func(t *testing.T) {
			cacheEnv(t)
			t.Setenv("ELESTIO_JWT_CACHE", "off")
			api, srv := newFakeAPI(t)
			api.rejectN = 100
			api.rejectRaw = fmt.Sprintf(`{"status":"KO","code":%q}`, code)

			client, _ := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
			_, err := client.Project.GetList()
			if err == nil || !strings.Contains(err.Error(), code) {
				t.Fatalf("the original 401 must reach the caller unchanged, got %v", err)
			}
			if n := api.signIn.calls.Load(); n != 1 {
				t.Fatalf("a %s response must not trigger a new sign-in (sign-ins=%d)", code, n)
			}
			if n := api.apiCalls.Load(); n != 1 {
				t.Fatalf("a %s response must not be replayed (calls=%d)", code, n)
			}
		})
	}
}

func TestAuthTransport_PersistentUnauthorizedDoesNotLoop(t *testing.T) {
	cacheEnv(t)
	t.Setenv("ELESTIO_JWT_CACHE", "off")
	api, srv := newFakeAPI(t)
	api.rejectN = 1000
	api.rejectRaw = `{"status":"KO","code":"Unauthorized"}`

	client, _ := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
	if _, err := client.Project.GetList(); err == nil {
		t.Fatal("expected an error")
	}
	if n := api.apiCalls.Load(); n != 2 {
		t.Fatalf("API calls = %d, want exactly 2 (original + one replay)", n)
	}
	if n := api.signIn.calls.Load(); n > 2 {
		t.Fatalf("sign-ins = %d, want at most 2", n)
	}
}

func TestNewAPIClient_UserSuppliedJWTNeedsNoSignIn(t *testing.T) {
	cacheEnv(t)
	api, srv := newFakeAPI(t)
	supplied := makeJWT(t, time.Now().Add(72*time.Hour), "supplied")

	client, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, jwt: supplied})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Project.GetList(); err != nil {
		t.Fatal(err)
	}
	if n := api.signIn.calls.Load(); n != 0 {
		t.Fatalf("sign-ins = %d, want 0", n)
	}
	if api.seen[0].header != supplied {
		t.Fatal("the supplied JWT was not used")
	}
}

func TestNewAPIClient_FailsAtConfigureWhenSignInIsRejected(t *testing.T) {
	cacheEnv(t)
	t.Setenv("ELESTIO_JWT_CACHE", "off")
	api, srv := newFakeAPI(t)
	api.signIn.status = 401
	api.signIn.body = `{"status":"KO","code":"InvalidCredentials","message":"The provided credentials are invalid."}`

	_, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
	if err == nil || !strings.Contains(err.Error(), "InvalidCredentials") {
		t.Fatalf("wrong credentials must fail up front with the API message, got %v", err)
	}
}

// A dropped connection produces a *url.Error that prints the request URL. The
// client's URL carries an empty token, so the real JWT can never appear there.
func TestAuthTransport_NetworkErrorsDoNotContainTheJWT(t *testing.T) {
	cacheEnv(t)
	t.Setenv("ELESTIO_JWT_CACHE", "off")
	api, srv := newFakeAPI(t)

	var drop atomic.Bool
	wrapped := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if drop.Load() && r.URL.Path != "/api/auth/checkAPIToken" {
			hj, _ := w.(http.Hijacker)
			c, _, _ := hj.Hijack()
			_ = c.Close()
			return
		}
		api.handler(w, r)
	}))
	defer wrapped.Close()
	_ = srv

	client, err := newAPIClient(context.Background(), clientConfig{baseURL: wrapped.URL, creds: testCreds})
	if err != nil {
		t.Fatal(err)
	}
	drop.Store(true)
	_, err = client.Project.GetList()
	if err == nil {
		t.Fatal("expected a network error")
	}
	token := api.signIn.token(1)
	if strings.Contains(err.Error(), token) || strings.Contains(err.Error(), "eyJ") {
		t.Fatalf("network error leaks the JWT: %v", err)
	}
}

func TestAuthTransport_NeverSendsTheTokenToAnotherHost(t *testing.T) {
	var got atomic.Value
	other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got.Store(r.Header.Get("Authorization") + "|" + r.URL.RawQuery)
	}))
	defer other.Close()

	f := &fakeSignIn{t: t}
	signIn := httptest.NewServer(http.HandlerFunc(f.handler))
	defer signIn.Close()
	src := newSource(t, signIn.URL, testCreds, nil)

	tr := &authTransport{base: http.DefaultTransport, host: "api.example.invalid", src: src}
	req, _ := http.NewRequest(http.MethodPost, other.URL+"/x?jwt=", strings.NewReader(`{"jwt":""}`))
	resp, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if v, _ := got.Load().(string); strings.Contains(v, "Bearer") || strings.Contains(v, "eyJ") {
		t.Fatalf("token sent to a foreign host: %q", v)
	}
	if f.calls.Load() != 0 {
		t.Fatal("a foreign-host request must not even trigger a sign-in")
	}
}

func TestInjectJWT(t *testing.T) {
	out := injectJWT([]byte(`{"projectID":"1","big":12345678901234567890,"jwt":""}`), "TOK")
	var m map[string]json.RawMessage
	if err := json.Unmarshal(out, &m); err != nil {
		t.Fatal(err)
	}
	if string(m["jwt"]) != `"TOK"` || string(m["projectID"]) != `"1"` || string(m["big"]) != `12345678901234567890` {
		t.Fatalf("unexpected body %s", out)
	}
	for _, in := range []string{`{"a":1}`, `not json`, `[1,2]`, ``} {
		if got := string(injectJWT([]byte(in), "TOK")); got != in {
			t.Errorf("%q must be left alone, got %q", in, got)
		}
	}
}
