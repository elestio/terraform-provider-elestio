package provider

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/elestio/elestio-go-api-client/v2"
)

const (
	// apiRequestTimeout bounds a single API round trip. The Elestio client
	// does not accept a context, so this is the only thing that stops a hung
	// connection from blocking Terraform (and the process) indefinitely.
	apiRequestTimeout = 2 * time.Minute

	// apiSignInTimeout bounds the initial authentication performed by
	// elestio.NewClient, which runs before a hardened HTTP client can be set.
	apiSignInTimeout = time.Minute

	maxAPIRedirects = 3
)

// newHardenedHTTPClient returns an HTTP client for the Elestio API with an
// overall timeout, TLS >= 1.2, bounded connection pooling, and a redirect
// policy that never leaves the original host or downgrades to plain HTTP.
// The API client authenticates with a JWT placed in the request URL, so a
// redirect to another host would disclose it.
func newHardenedHTTPClient() *http.Client {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: apiRequestTimeout,
		ExpectContinueTimeout: time.Second,
		MaxIdleConns:          20,
		MaxIdleConnsPerHost:   10,
		IdleConnTimeout:       90 * time.Second,
	}

	return &http.Client{
		Timeout:       apiRequestTimeout,
		Transport:     transport,
		CheckRedirect: checkRedirect,
	}
}

func checkRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= maxAPIRedirects {
		return fmt.Errorf("stopped after %d redirects", maxAPIRedirects)
	}
	first := via[0].URL
	if req.URL.Scheme != "https" && first.Scheme == "https" {
		return errors.New("refusing redirect from https to a non-https URL")
	}
	if req.URL.Host != first.Host {
		return fmt.Errorf("refusing redirect to a different host (%s)", req.URL.Host)
	}
	return nil
}

// clientConfig describes how to reach and authenticate to the Elestio API.
type clientConfig struct {
	baseURL string // empty selects the production API
	creds   apiCredentials
	jwt     string // optional pre-obtained session JWT (ELESTIO_JWT)
}

// newAPIClient returns an API client that reuses one session JWT across
// requests and across Terraform runs instead of signing in every time (the
// backend allows only 15 sign-ins per hour). It verifies authentication up
// front so a wrong token fails at provider configuration.
func newAPIClient(ctx context.Context, cfg clientConfig) (*elestio.Client, error) {
	if cfg.baseURL == "" {
		cfg.baseURL = elestio.BaseURLV1
	}
	base, err := url.Parse(cfg.baseURL)
	if err != nil || base.Host == "" {
		return nil, fmt.Errorf("invalid API base URL")
	}

	signInClient := newHardenedHTTPClient()
	signInClient.Timeout = apiSignInTimeout
	src := newTokenSource(cfg.baseURL, cfg.creds, signInClient, newJWTCache(cfg.creds))
	if cfg.jwt != "" {
		src.seed(cfg.jwt)
	}
	if _, err := src.Token(ctx); err != nil {
		return nil, err
	}

	httpClient := newHardenedHTTPClient()
	httpClient.Transport = &authTransport{base: httpClient.Transport, host: base.Host, src: src}

	// NewUnsignedClient performs no sign-in; its requests carry an empty JWT
	// that authTransport replaces with the real one.
	client := elestio.NewUnsignedClient()
	client.BaseURL = cfg.baseURL
	client.HTTPClient = httpClient
	return client, nil
}
