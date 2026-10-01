package utils

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"testing"
)

func TestRedactSecrets(t *testing.T) {
	const token = "eyJhbGciOiJIUzI1NiJ9.payload.sig-nature_1"

	tests := map[string]string{
		"url error":               (&url.Error{Op: "Post", URL: "https://api.elest.io/api/x?jwt=" + token, Err: errors.New("context deadline exceeded")}).Error(),
		"url error, other params": "Get \"https://api.elest.io/x?a=1&jwt=" + token + "&b=2\": EOF",
		"json body":               `{"jwt":"` + token + `","status":"OK"}`,
		"bearer header":           "Authorization: Bearer " + token,
	}

	for name, in := range tests {
		t.Run(name, func(t *testing.T) {
			out := RedactSecrets(in)
			if strings.Contains(out, token) {
				t.Fatalf("token leaked: %q", out)
			}
			if !strings.Contains(out, "REDACTED") {
				t.Fatalf("expected marker in %q", out)
			}
		})
	}

	if got := RedactSecrets("Get \"https://api.elest.io/x?a=1&b=2\": EOF"); strings.Contains(got, "REDACTED") {
		t.Fatalf("unexpected redaction: %q", got)
	}
}

func TestRedactError(t *testing.T) {
	if RedactError(nil) != "" {
		t.Fatal("nil error must give empty string")
	}
	err := fmt.Errorf("wrapped: %w", &url.Error{Op: "Get", URL: "https://h/?jwt=secret", Err: errors.New("boom")})
	if got := RedactError(err); strings.Contains(got, "secret") {
		t.Fatalf("leaked: %q", got)
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		err  error
		want bool
	}{
		{nil, false},
		{errors.New("service not found"), true},
		{errors.New("project not found"), true},
		{errors.New("load balancer not found"), true},
		{fmt.Errorf("ctx: %w", errors.New("service not found")), true},
		{errors.New("request failed with status code 500: service not found"), false},
		{errors.New("dial tcp: connection refused"), false},
		{errors.New("request failed with status code 401: unauthorized"), false},
		// Real backend responses (api/servers/getServerDetails.js).
		{errors.New(`request failed with status code 401: {"status":"KO","code":"InvalidServer","message":"Invalid Server; unable to retrieve information."}`), true},
		{errors.New(`request failed with status code 401: {"status":"KO","code":"service_deleted","message":"Sevice deleting"}`), true},
		{fmt.Errorf("get: %w", errors.New(`request failed with status code 401: {"status": "KO", "code": "service_deleted"}`)), true},
		// Other 401s are real failures, not proof of deletion.
		{errors.New(`request failed with status code 401: {"status":"KO","code":"Unauthorized","message":"Authentication failed."}`), false},
		{errors.New(`request failed with status code 401: {"status":"KO","code":"InvalidPermission"}`), false},
		{errors.New(`request failed with status code 401: {"status":"KO","code":"MissingPermission"}`), false},
		{errors.New(`request failed with status code 429: {"code":"TooManyRequests"}`), false},
		{errors.New(`request failed with status code 500: {"code":"InvalidServer"}`), false},
	}
	for _, tt := range tests {
		if got := IsNotFound(tt.err); got != tt.want {
			t.Errorf("IsNotFound(%v) = %v, want %v", tt.err, got, tt.want)
		}
	}
}
