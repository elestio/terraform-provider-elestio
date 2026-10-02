package provider

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"
)

const maxAuthBody = 8 << 20

var unauthorizedCode = regexp.MustCompile(`"code"\s*:\s*"Unauthorized"`)

// authTransport gives every API request the current session JWT.
//
// The Elestio API client builds each request with the JWT in the JSON body, in
// a `jwt` query parameter and in the Authorization header, and takes it from a
// private field set only by its own sign-in. This provider signs in itself (so
// the token can be reused), builds the client without a JWT, and lets this
// transport fill in the real one. Because the client's own request URL carries
// an empty token, Go's network errors can no longer print the JWT.
//
// If the API answers 401 with code Unauthorized (expired or revoked token) the
// request is replayed once with a fresh token. That is safe for non-idempotent
// calls too: the backend rejects the token before doing anything.
type authTransport struct {
	base http.RoundTripper
	host string
	src  *tokenSource
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Never attach the token to a request for another host.
	if req.URL.Host != t.host {
		return t.base.RoundTrip(req)
	}

	body, err := readRequestBody(req)
	if err != nil {
		return nil, err
	}

	token, err := t.src.Token(req.Context())
	if err != nil {
		return nil, err
	}
	resp, err := t.send(req, body, token)
	if err != nil || resp.StatusCode != http.StatusUnauthorized {
		return resp, err
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
	_ = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if unauthorizedCode.Match(respBody) {
		if fresh, rerr := t.src.Refresh(req.Context(), token); rerr == nil && fresh != token {
			return t.send(req, body, fresh)
		}
	}
	resp.Body = io.NopCloser(bytes.NewReader(respBody))
	return resp, nil
}

func (t *authTransport) send(req *http.Request, body []byte, token string) (*http.Response, error) {
	clone := req.Clone(req.Context())
	clone.Header.Set("Authorization", "Bearer "+token)
	q := clone.URL.Query()
	q.Set("jwt", token)
	clone.URL.RawQuery = q.Encode()

	if body != nil {
		rewritten := injectJWT(body, token)
		clone.Body = io.NopCloser(bytes.NewReader(rewritten))
		clone.ContentLength = int64(len(rewritten))
		clone.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(rewritten)), nil }
		clone.Header.Set("Content-Length", strconv.Itoa(len(rewritten)))
	}
	return t.base.RoundTrip(clone)
}

func readRequestBody(req *http.Request) ([]byte, error) {
	if req.Body == nil || req.Body == http.NoBody {
		return nil, nil
	}
	defer func() { _ = req.Body.Close() }()
	b, err := io.ReadAll(io.LimitReader(req.Body, maxAuthBody))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// injectJWT sets the "jwt" field of a JSON object body, but only if the client
// put one there. Anything else is returned unchanged.
func injectJWT(body []byte, token string) []byte {
	var fields map[string]json.RawMessage
	if json.Unmarshal(body, &fields) != nil {
		return body
	}
	if _, ok := fields["jwt"]; !ok {
		return body
	}
	quoted, err := json.Marshal(token)
	if err != nil {
		return body
	}
	fields["jwt"] = quoted
	out, err := json.Marshal(fields)
	if err != nil {
		return body
	}
	return out
}
