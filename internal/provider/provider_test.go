package provider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func TestProviderSchemaMarksAPITokenSensitive(t *testing.T) {
	var resp provider.SchemaResponse
	New("test")().Schema(context.Background(), provider.SchemaRequest{}, &resp)

	attr, ok := resp.Schema.Attributes["api_token"]
	if !ok || !attr.IsSensitive() {
		t.Fatal("api_token must exist and be sensitive")
	}
}

// Every resource must have a valid schema, a unique type name, and no secret
// attribute that Terraform would print in plans.
func TestAllResourceSchemas(t *testing.T) {
	ctx := context.Background()
	resources := New("test")().Resources(ctx)
	if len(resources) == 0 {
		t.Fatal("no resources registered")
	}

	seen := map[string]bool{}
	for _, newResource := range resources {
		r := newResource()

		var meta resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "elestio"}, &meta)
		if seen[meta.TypeName] {
			// Deprecated aliases share template IDs but never type names.
			t.Errorf("duplicate resource type name %q", meta.TypeName)
		}
		seen[meta.TypeName] = true

		var sr resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &sr)
		if sr.Diagnostics.HasError() {
			t.Errorf("%s: schema diagnostics: %v", meta.TypeName, sr.Diagnostics)
			continue
		}
		if diags := sr.Schema.ValidateImplementation(ctx); diags.HasError() {
			t.Errorf("%s: invalid schema implementation: %v", meta.TypeName, diags)
		}

		assertSecretsSensitive(t, meta.TypeName, sr.Schema.Attributes, false)
	}
}

func assertSecretsSensitive(t *testing.T, prefix string, attrs map[string]rschema.Attribute, parentSensitive bool) {
	t.Helper()
	for name, a := range attrs {
		full := prefix + "." + name
		sensitive := parentSensitive || a.IsSensitive()

		lower := strings.ToLower(name)
		if (strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "token")) && !sensitive {
			t.Errorf("%s holds a secret but is not sensitive", full)
		}

		if nested, ok := a.(rschema.SingleNestedAttribute); ok {
			assertSecretsSensitive(t, full, nested.Attributes, sensitive)
		}
	}
}
