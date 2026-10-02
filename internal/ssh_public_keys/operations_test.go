package ssh_public_keys

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type apiCall struct{ path, body string }

// fakeKeysAPI records requests. failAddOnce makes the first add-key call of
// the given key data fail with a 500.
func fakeKeysAPI(t *testing.T, failKeyData string) (*elestio.Client, func() []apiCall) {
	t.Helper()
	var mu sync.Mutex
	var calls []apiCall
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		calls = append(calls, apiCall{r.URL.Path, string(b)})
		mu.Unlock()
		if failKeyData != "" && strings.Contains(string(b), failKeyData) {
			http.Error(w, "boom", http.StatusInternalServerError)
			return
		}
		_, _ = w.Write([]byte(`{"status":"OK"}`))
	}))
	t.Cleanup(srv.Close)
	c := elestio.NewUnsignedClient()
	c.BaseURL = srv.URL
	c.HTTPClient = srv.Client()
	return c, func() []apiCall {
		mu.Lock()
		defer mu.Unlock()
		return append([]apiCall(nil), calls...)
	}
}

func key(user, data string) models.SSHPublicKeyModel {
	return models.SSHPublicKeyModel{Username: types.StringValue(user), KeyData: types.StringValue(data)}
}

func TestApplyChanges_FailedReplaceRestoresPreviousKey(t *testing.T) {
	client, calls := fakeKeysAPI(t, "ssh-ed25519 NEWKEY")

	err := ApplyChanges(context.Background(), "svc-1",
		nil, []models.SSHPublicKeyModel{key("alice", "ssh-ed25519 NEWKEY")}, nil,
		[]models.SSHPublicKeyModel{key("alice", "ssh-ed25519 OLDKEY")},
		"hetzner", client, nil, &elestio.Service{}, time.Second)

	if err == nil || !strings.Contains(err.Error(), "previous key was restored") {
		t.Fatalf("expected restore message, got %v", err)
	}

	// remove alice, add NEWKEY (fails), add OLDKEY (restore)
	var adds []string
	for _, c := range calls() {
		var m map[string]any
		_ = json.Unmarshal([]byte(c.body), &m)
		if k, ok := m["key"].(string); ok {
			adds = append(adds, k)
		}
	}
	if len(adds) < 2 || !strings.Contains(adds[len(adds)-1], "OLDKEY") {
		t.Fatalf("old key must be re-added last, adds=%v (calls=%v)", adds, calls())
	}
}

func TestApplyChanges_NoChangesMakesNoCalls(t *testing.T) {
	client, calls := fakeKeysAPI(t, "")
	err := ApplyChanges(context.Background(), "svc-1", nil, nil, nil, nil, "scaleway", client, nil, &elestio.Service{}, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if n := len(calls()); n != 0 {
		t.Fatalf("expected no API calls (and no scaleway reboot), got %d", n)
	}
}
