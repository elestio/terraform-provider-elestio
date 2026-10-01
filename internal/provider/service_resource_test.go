package provider

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// fakeAPI starts an HTTP server standing in for the Elestio API and returns a
// client pointed at it. handler receives every request path.
func fakeAPI(t *testing.T, handler func(w http.ResponseWriter, r *http.Request)) *elestio.Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(srv.Close)

	c := elestio.NewUnsignedClient()
	c.BaseURL = srv.URL
	c.HTTPClient = srv.Client()
	return c
}

// serviceNotFoundHandler mimics the real backend for a deleted or unknown
// service: HTTP 401 with code InvalidServer (api/servers/getServerDetails.js).
func serviceNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"status":"KO","code":"InvalidServer","message":"Invalid Server; unable to retrieve information."}`))
}

// unauthorizedHandler mimics an expired or invalid session token.
func unauthorizedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"status":"KO","code":"Unauthorized","message":"Authentication failed."}`))
}

func serverErrorHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "boom", http.StatusInternalServerError)
}

func newTestServiceResource(t *testing.T, c *elestio.Client) (*ServiceResource, resource.SchemaResponse) {
	t.Helper()
	r := NewServiceResource(&ServiceTemplate{TemplateId: 1, ResourceName: "test", DocumentationName: "Test", DefaultVersion: "latest"}).(*ServiceResource)
	r.client = c
	r.deletionPollDelay = time.Millisecond
	r.deletionPollMin = time.Millisecond

	var sr resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &sr)
	if sr.Diagnostics.HasError() {
		t.Fatalf("schema: %v", sr.Diagnostics)
	}
	return r, sr
}

// stateWithIDs builds a state in which every attribute is null except the
// service and project identifiers.
func stateWithIDs(t *testing.T, sr resource.SchemaResponse) tfsdk.State {
	t.Helper()
	ctx := context.Background()
	objType := sr.Schema.Type().TerraformType(ctx).(tftypes.Object)
	vals := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	st := tfsdk.State{Schema: sr.Schema, Raw: tftypes.NewValue(objType, vals)}
	if d := st.SetAttribute(ctx, path.Root("id"), "svc-1"); d.HasError() {
		t.Fatal(d)
	}
	if d := st.SetAttribute(ctx, path.Root("project_id"), "prj-1"); d.HasError() {
		t.Fatal(d)
	}
	return st
}

func TestServiceRead_MissingServiceIsRemovedFromState(t *testing.T) {
	r, sr := newTestServiceResource(t, fakeAPI(t, serviceNotFoundHandler))
	st := stateWithIDs(t, sr)

	resp := resource.ReadResponse{State: st}
	r.Read(context.Background(), resource.ReadRequest{State: st}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("a deleted service must not be an error: %v", resp.Diagnostics)
	}
	if !resp.State.Raw.IsNull() {
		t.Fatal("state must be removed so Terraform plans a re-create")
	}
}

func TestServiceRead_ServerErrorKeepsStateAndFails(t *testing.T) {
	r, sr := newTestServiceResource(t, fakeAPI(t, serverErrorHandler))
	st := stateWithIDs(t, sr)

	resp := resource.ReadResponse{State: st}
	r.Read(context.Background(), resource.ReadRequest{State: st}, &resp)

	if !resp.Diagnostics.HasError() {
		t.Fatal("a 500 must surface as an error")
	}
	if resp.State.Raw.IsNull() {
		t.Fatal("a transient API failure must never drop the resource from state")
	}
}

func TestServiceRead_ExpiredTokenKeepsState(t *testing.T) {
	r, sr := newTestServiceResource(t, fakeAPI(t, unauthorizedHandler))
	st := stateWithIDs(t, sr)

	resp := resource.ReadResponse{State: st}
	r.Read(context.Background(), resource.ReadRequest{State: st}, &resp)
	if !resp.Diagnostics.HasError() || resp.State.Raw.IsNull() {
		t.Fatal("an auth failure must error and must not remove the resource from state")
	}
}

func TestServiceDelete_AlreadyDeletedSucceeds(t *testing.T) {
	r, sr := newTestServiceResource(t, fakeAPI(t, serviceNotFoundHandler))
	st := stateWithIDs(t, sr)

	var resp resource.DeleteResponse
	r.Delete(context.Background(), resource.DeleteRequest{State: st}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("deleting an already-deleted service must succeed: %v", resp.Diagnostics)
	}
}

func TestWaitServiceDeletion(t *testing.T) {
	svc := &elestio.Service{ProjectID: "prj-1", ID: "svc-1"}

	t.Run("not found means deleted", func(t *testing.T) {
		r, _ := newTestServiceResource(t, fakeAPI(t, serviceNotFoundHandler))
		if err := r.waitServiceDeletion(context.Background(), svc, 5*time.Second); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("server errors are not proof of deletion", func(t *testing.T) {
		r, _ := newTestServiceResource(t, fakeAPI(t, serverErrorHandler))
		if err := r.waitServiceDeletion(context.Background(), svc, 300*time.Millisecond); err == nil {
			t.Fatal("a 500 response must not be treated as a completed deletion")
		}
	})

	t.Run("an expired token (401 Unauthorized) is not proof of deletion", func(t *testing.T) {
		r, _ := newTestServiceResource(t, fakeAPI(t, unauthorizedHandler))
		if err := r.waitServiceDeletion(context.Background(), svc, 300*time.Millisecond); err == nil {
			t.Fatal("401 Unauthorized must not be treated as a completed deletion")
		}
	})

	t.Run("context cancellation stops the wait", func(t *testing.T) {
		r, _ := newTestServiceResource(t, fakeAPI(t, serverErrorHandler))
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		if err := r.waitServiceDeletion(ctx, svc, time.Minute); err == nil {
			t.Fatal("expected cancellation error")
		}
	})
}

// Creating a server is not idempotent and billed per attempt: a failed
// request must never be re-sent.
func TestServiceCreate_IsNotRetried(t *testing.T) {
	var calls atomic.Int32
	c := fakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/createServer") {
			calls.Add(1)
		}
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	r, sr := newTestServiceResource(t, c)

	plan := tfsdk.Plan{Schema: sr.Schema, Raw: stateWithIDs(t, sr).Raw}
	resp := resource.CreateResponse{State: tfsdk.State{Schema: sr.Schema, Raw: stateWithIDs(t, sr).Raw}}

	done := make(chan struct{})
	go func() {
		defer close(done)
		r.Create(context.Background(), resource.CreateRequest{Plan: plan}, &resp)
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("Create kept retrying")
	}

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error")
	}
	if n := calls.Load(); n != 1 {
		t.Fatalf("createServer called %d times, want exactly 1", n)
	}
}

func TestProjectRead_MissingProjectIsRemovedFromState(t *testing.T) {
	c := fakeAPI(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"OK","data":{"projects":[]}}`))
	})
	r := &ProjectResource{client: c}

	var sr resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &sr)
	ctx := context.Background()
	objType := sr.Schema.Type().TerraformType(ctx).(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	st := tfsdk.State{Schema: sr.Schema, Raw: tftypes.NewValue(objType, vals)}
	if d := st.SetAttribute(ctx, path.Root("id"), "prj-1"); d.HasError() {
		t.Fatal(d)
	}

	resp := resource.ReadResponse{State: st}
	r.Read(ctx, resource.ReadRequest{State: st}, &resp)
	if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() {
		t.Fatalf("missing project must be removed from state without error (diags: %v)", resp.Diagnostics)
	}
}
