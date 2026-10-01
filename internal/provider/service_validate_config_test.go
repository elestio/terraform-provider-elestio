package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestServiceValidateConfig_Backups(t *testing.T) {
	r, sr := newTestServiceResource(t, nil)
	ctx := context.Background()
	objType := sr.Schema.Type().TerraformType(ctx).(tftypes.Object)

	cfg := func(backups, level any) tfsdk.Config {
		vals := map[string]tftypes.Value{}
		for name, typ := range objType.AttributeTypes {
			vals[name] = tftypes.NewValue(typ, nil)
		}
		if backups != nil {
			vals["backups_enabled"] = tftypes.NewValue(tftypes.Bool, backups)
		}
		if level != nil {
			vals["support_level"] = tftypes.NewValue(tftypes.String, level)
		}
		return tfsdk.Config{Schema: sr.Schema, Raw: tftypes.NewValue(objType, vals)}
	}

	tests := []struct {
		name    string
		config  tfsdk.Config
		wantErr bool
	}{
		{"backups with default level1 (support_level unset)", cfg(true, nil), true},
		{"backups with explicit level1", cfg(true, "level1"), true},
		{"backups with level2", cfg(true, "level2"), false},
		{"backups with unknown level defers to plan", cfg(true, tftypes.UnknownValue), false},
		{"no backups", cfg(false, nil), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp resource.ValidateConfigResponse
			r.ValidateConfig(ctx, resource.ValidateConfigRequest{Config: tt.config}, &resp)
			if got := resp.Diagnostics.HasError(); got != tt.wantErr {
				t.Fatalf("error=%v, want %v: %v", got, tt.wantErr, resp.Diagnostics)
			}
		})
	}
}
