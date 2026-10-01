package provider

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
)

func TestHardenedHTTPClientSettings(t *testing.T) {
	c := newHardenedHTTPClient()

	if c.Timeout <= 0 {
		t.Fatal("client must have an overall timeout")
	}
	tr, ok := c.Transport.(*http.Transport)
	if !ok {
		t.Fatalf("unexpected transport %T", c.Transport)
	}
	if tr.TLSClientConfig == nil || tr.TLSClientConfig.MinVersion < tls.VersionTLS12 {
		t.Fatal("TLS 1.2 must be the minimum version")
	}
	if tr.TLSClientConfig.InsecureSkipVerify {
		t.Fatal("certificate verification must stay enabled")
	}
	if tr.TLSHandshakeTimeout <= 0 || tr.IdleConnTimeout <= 0 || tr.MaxIdleConns <= 0 {
		t.Fatal("transport must bound handshake and idle connections")
	}

	// Two clients must not share mutable state.
	if newHardenedHTTPClient() == c || newHardenedHTTPClient().Transport == c.Transport {
		t.Fatal("clients must not share a transport")
	}
}

func TestHardenedHTTPClientRedirects(t *testing.T) {
	var otherHits atomic.Int32
	other := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		otherHits.Add(1)
	}))
	defer other.Close()

	var loops atomic.Int32
	var srv *httptest.Server
	srv = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/cross":
			http.Redirect(w, r, other.URL+"/?jwt=secret", http.StatusTemporaryRedirect)
		case "/loop":
			loops.Add(1)
			http.Redirect(w, r, srv.URL+"/loop", http.StatusFound)
		case "/same":
			http.Redirect(w, r, srv.URL+"/ok", http.StatusFound)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	client := newHardenedHTTPClient()
	// Trust the httptest certificate (shared by both servers) in this copy only.
	pool := x509.NewCertPool()
	pool.AddCert(srv.Certificate())
	pool.AddCert(other.Certificate())
	client.Transport.(*http.Transport).TLSClientConfig.RootCAs = pool

	t.Run("cross host redirect is refused", func(t *testing.T) {
		_, err := client.Get(srv.URL + "/cross?jwt=secret")
		if err == nil || !strings.Contains(err.Error(), "different host") {
			t.Fatalf("expected refusal, got %v", err)
		}
		if otherHits.Load() != 0 {
			t.Fatal("request reached the redirect target")
		}
	})

	t.Run("redirect loop is bounded", func(t *testing.T) {
		_, err := client.Get(srv.URL + "/loop")
		if err == nil || !strings.Contains(err.Error(), "redirects") {
			t.Fatalf("expected redirect limit error, got %v", err)
		}
		if n := loops.Load(); n > maxAPIRedirects+1 {
			t.Fatalf("followed %d redirects", n)
		}
	})

	t.Run("same host redirect is followed", func(t *testing.T) {
		resp, err := client.Get(srv.URL + "/same")
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = resp.Body.Close() }()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
	})
}

func TestCheckRedirectRefusesDowngrade(t *testing.T) {
	first, _ := http.NewRequest(http.MethodGet, "https://api.example.com/a", nil)
	next, _ := http.NewRequest(http.MethodGet, "http://api.example.com/b", nil)
	if err := checkRedirect(next, []*http.Request{first}); err == nil {
		t.Fatal("https -> http redirect must be refused")
	}
}
