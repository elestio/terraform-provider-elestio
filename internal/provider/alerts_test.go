package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// backendRequiredRules are the parameters the backend's enableMonicAlert reads
// unconditionally (api helper/users/doActionOnServer.js). A missing one makes
// the backend throw and answer 500 ActionError.
var backendRequiredRules = []string{
	"CPU", "MEMORY", "SWAP", "SPACE", "INODE", "READ_RATE", "READ_RATE_OPERATIONS",
	"WRITE_RATE", "WRITE_RATE_OPERATIONS", "SATURATION", "DOWNLOAD", "UPLOAD",
}

// fakeBackendAlerts behaves like the real DoActionOnServer enableAlerts branch.
func fakeBackendAlerts(t *testing.T, sawJWT *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/auth/checkAPIToken" {
			_, _ = w.Write([]byte(`{"status":"OK","jwt":"` + makeJWT(t, time.Now().Add(time.Hour), "a") + `"}`))
			return
		}
		b, _ := io.ReadAll(r.Body)
		var req struct {
			JWT                 string `json:"jwt"`
			Action              string `json:"action"`
			VMID                string `json:"vmID"`
			Rules               string `json:"rules"`
			MonitCycleInSeconds int    `json:"monitCycleInSeconds"`
		}
		_ = json.Unmarshal(b, &req)
		*sawJWT = req.JWT

		var rules []alertRule
		if err := json.Unmarshal([]byte(req.Rules), &rules); err != nil || req.Action != "enableAlerts" || req.MonitCycleInSeconds <= 0 {
			w.WriteHeader(500)
			_, _ = w.Write([]byte(`{"status":"KO","code":"ActionError","message":"An error occured, please try again in few minutes ..."}`))
			return
		}
		byName := map[string]alertRule{}
		for _, ru := range rules {
			byName[ru.Parameter] = ru
		}
		for _, need := range backendRequiredRules {
			if _, ok := byName[need]; !ok { // ruleResult.X.value would throw
				w.WriteHeader(500)
				_, _ = w.Write([]byte(`{"status":"KO","code":"ActionError","message":"An error occured, please try again in few minutes ..."}`))
				return
			}
		}
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	}
}

func TestEnableAlerts_SendsEveryRuleTheBackendReads(t *testing.T) {
	cacheEnv(t)
	var sawJWT string
	srv := httptest.NewServer(fakeBackendAlerts(t, &sawJWT))
	defer srv.Close()

	client, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
	if err != nil {
		t.Fatal(err)
	}
	if err := enableAlerts(context.Background(), client, "1234567890"); err != nil {
		t.Fatalf("enableAlerts must succeed against a backend that requires all 12 rules: %v", err)
	}
	if sawJWT == "" || strings.Count(sawJWT, ".") != 2 {
		t.Fatalf("the transport must put the real JWT in the body, got %q", sawJWT)
	}
}

// Regression for the live failure: the API client's own EnableAlerts sends ten
// rules and the backend answers 500 ActionError.
func TestAPIClientEnableAlertsAloneFailsAgainstCurrentBackend(t *testing.T) {
	cacheEnv(t)
	var sawJWT string
	srv := httptest.NewServer(fakeBackendAlerts(t, &sawJWT))
	defer srv.Close()

	client, err := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
	if err != nil {
		t.Fatal(err)
	}
	if err := client.Service.EnableAlerts("1234567890"); err == nil {
		t.Skip("elestio-go-api-client now sends all required rules; the workaround in alerts.go can be removed")
	}
}

func TestDefaultAlertRules_MatchBackendDefaults(t *testing.T) {
	got := map[string]bool{}
	for _, r := range defaultAlertRules {
		got[r.Parameter] = true
		if r.Value <= 0 || r.Cycles <= 0 || r.Unit == "" {
			t.Errorf("rule %s has an empty field: %+v", r.Parameter, r)
		}
	}
	for _, need := range backendRequiredRules {
		if !got[need] {
			t.Errorf("missing rule %s", need)
		}
	}
	if len(defaultAlertRules) != len(backendRequiredRules) {
		t.Errorf("have %d rules, backend defines %d", len(defaultAlertRules), len(backendRequiredRules))
	}
}

func TestEnableAlerts_SurfacesBackendFailures(t *testing.T) {
	cacheEnv(t)
	for name, tc := range map[string]struct {
		status int
		body   string
	}{
		"500":         {500, `{"status":"KO","code":"ActionError","message":"boom"}`},
		"KO with 200": {200, `{"status":"KO","message":"not allowed"}`},
	} {
		t.Run(name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/api/auth/checkAPIToken" {
					_, _ = w.Write([]byte(`{"status":"OK","jwt":"` + makeJWT(t, time.Now().Add(time.Hour), "a") + `"}`))
					return
				}
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()
			client, _ := newAPIClient(context.Background(), clientConfig{baseURL: srv.URL, creds: testCreds})
			if err := enableAlerts(context.Background(), client, "1"); err == nil {
				t.Fatal("expected an error")
			}
		})
	}
}
