package validators

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsEmail(t *testing.T) {
	tests := map[string]bool{
		"user@example.com":                    true,
		"first.last+tag@sub.example.io":       true,
		"":                                    false,
		"not-an-email":                        false,
		"Name <user@example.com>":             false,
		"user@example.com, other@example.com": false,
		"user@example.com\r\nBcc: x@y.z":      false,
		"user@@example.com":                   false,
	}
	for in, valid := range tests {
		req := validator.StringRequest{Path: path.Root("e"), ConfigValue: types.StringValue(in)}
		var resp validator.StringResponse
		IsEmail().ValidateString(context.Background(), req, &resp)
		if got := !resp.Diagnostics.HasError(); got != valid {
			t.Errorf("%q: valid=%v, want %v", in, got, valid)
		}
	}

	for _, v := range []types.String{types.StringNull(), types.StringUnknown()} {
		var resp validator.StringResponse
		IsEmail().ValidateString(context.Background(), validator.StringRequest{ConfigValue: v}, &resp)
		if resp.Diagnostics.HasError() {
			t.Errorf("null/unknown must be skipped")
		}
	}
}
