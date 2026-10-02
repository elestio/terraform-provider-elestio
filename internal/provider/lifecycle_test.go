package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func emptyState(t *testing.T, sr resource.SchemaResponse) tfsdk.State {
	t.Helper()
	objType := sr.Schema.Type().TerraformType(context.Background()).(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	return tfsdk.State{Schema: sr.Schema, Raw: tftypes.NewValue(objType, vals)}
}

func TestImportState_RejectsMalformedIDs(t *testing.T) {
	ctx := context.Background()
	svc, svcSchema := newTestServiceResource(t, nil)
	lb := &LoadBalancerResource{}
	var lbSchema resource.SchemaResponse
	lb.Schema(ctx, resource.SchemaRequest{}, &lbSchema)

	importers := map[string]struct {
		imp    resource.ResourceWithImportState
		schema resource.SchemaResponse
	}{
		"service":       {svc, svcSchema},
		"load balancer": {lb, lbSchema},
	}

	bad := []string{"", "onlyone", ",svc", "prj,", "a,b,c", ",", "a b"[:1]}
	for name, im := range importers {
		for _, id := range bad {
			var resp resource.ImportStateResponse
			resp.State = emptyState(t, im.schema)
			im.imp.ImportState(ctx, resource.ImportStateRequest{ID: id}, &resp)
			if !resp.Diagnostics.HasError() {
				t.Errorf("%s: import id %q must be rejected", name, id)
			}
		}

		var resp resource.ImportStateResponse
		resp.State = emptyState(t, im.schema)
		im.imp.ImportState(ctx, resource.ImportStateRequest{ID: "prj-1,svc-1"}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("%s: valid id rejected: %v", name, resp.Diagnostics)
		}
		var pid, id string
		resp.State.GetAttribute(ctx, path.Root("project_id"), &pid)
		resp.State.GetAttribute(ctx, path.Root("id"), &id)
		if pid != "prj-1" || id != "svc-1" {
			t.Errorf("%s: got project_id=%q id=%q", name, pid, id)
		}
	}
}

func TestLoadBalancerRead_MissingIsRemovedFromState(t *testing.T) {
	ctx := context.Background()
	lb := &LoadBalancerResource{client: fakeAPI(t, serviceNotFoundHandler)}
	var sr resource.SchemaResponse
	lb.Schema(ctx, resource.SchemaRequest{}, &sr)
	st := emptyState(t, sr)
	// ConfigModel is a non-pointer field, so "config" must be a non-null object.
	cfgType := sr.Schema.Type().TerraformType(ctx).(tftypes.Object).AttributeTypes["config"].(tftypes.Object)
	cfgVals := map[string]tftypes.Value{}
	for name, typ := range cfgType.AttributeTypes {
		cfgVals[name] = tftypes.NewValue(typ, nil)
	}
	root := st.Raw.Copy()
	var rootVals map[string]tftypes.Value
	_ = root.As(&rootVals)
	rootVals["config"] = tftypes.NewValue(cfgType, cfgVals)
	st.Raw = tftypes.NewValue(root.Type(), rootVals)
	_ = st.SetAttribute(ctx, path.Root("id"), "lb-1")
	_ = st.SetAttribute(ctx, path.Root("project_id"), "prj-1")

	resp := resource.ReadResponse{State: st}
	lb.Read(ctx, resource.ReadRequest{State: st}, &resp)
	if resp.Diagnostics.HasError() || !resp.State.Raw.IsNull() {
		t.Fatalf("missing load balancer must be removed from state: %v", resp.Diagnostics)
	}

	lb500 := &LoadBalancerResource{client: fakeAPI(t, serverErrorHandler)}
	resp = resource.ReadResponse{State: st}
	lb500.Read(ctx, resource.ReadRequest{State: st}, &resp)
	if !resp.Diagnostics.HasError() || resp.State.Raw.IsNull() {
		t.Fatal("a 500 must fail the read and keep state")
	}
}

func TestProviderConfigure_MissingCredentials(t *testing.T) {
	t.Setenv("ELESTIO_EMAIL", "")
	t.Setenv("ELESTIO_API_TOKEN", "")

	ctx := context.Background()
	p := New("test")()
	var sr provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &sr)

	objType := sr.Schema.Type().TerraformType(ctx).(tftypes.Object)
	vals := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		vals[name] = tftypes.NewValue(typ, nil)
	}
	cfg := tfsdk.Config{Schema: sr.Schema, Raw: tftypes.NewValue(objType, vals)}

	var resp provider.ConfigureResponse
	p.Configure(ctx, provider.ConfigureRequest{Config: cfg}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("missing credentials must be an error")
	}
	if resp.ResourceData != nil {
		t.Fatal("no client may be handed to resources without credentials")
	}
}

func TestProviderConfigure_UnknownCredentials(t *testing.T) {
	ctx := context.Background()
	p := New("test")()
	var sr provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &sr)

	objType := sr.Schema.Type().TerraformType(ctx).(tftypes.Object)
	cfg := tfsdk.Config{Schema: sr.Schema, Raw: tftypes.NewValue(objType, map[string]tftypes.Value{
		"email":     tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"api_token": tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
	})}
	var resp provider.ConfigureResponse
	p.Configure(ctx, provider.ConfigureRequest{Config: cfg}, &resp)
	if !resp.Diagnostics.HasError() {
		t.Fatal("unknown credentials must be an error")
	}
}
