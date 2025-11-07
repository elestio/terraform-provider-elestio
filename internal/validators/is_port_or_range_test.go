package validators

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsPortOrRange(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input               types.String
		expectedDiagnostics diag.Diagnostics
	}{
		"unknown": {
			input:               types.StringUnknown(),
			expectedDiagnostics: nil,
		},
		"null": {
			input:               types.StringNull(),
			expectedDiagnostics: nil,
		},
		// Valid single ports
		"valid-port-80": {
			input:               types.StringValue("80"),
			expectedDiagnostics: nil,
		},
		"valid-port-443": {
			input:               types.StringValue("443"),
			expectedDiagnostics: nil,
		},
		"valid-port-22": {
			input:               types.StringValue("22"),
			expectedDiagnostics: nil,
		},
		"valid-port-8080": {
			input:               types.StringValue("8080"),
			expectedDiagnostics: nil,
		},
		"valid-port-1": {
			input:               types.StringValue("1"),
			expectedDiagnostics: nil,
		},
		"valid-port-65535": {
			input:               types.StringValue("65535"),
			expectedDiagnostics: nil,
		},
		// Valid port ranges
		"valid-range-8000-9000": {
			input:               types.StringValue("8000-9000"),
			expectedDiagnostics: nil,
		},
		"valid-range-1-65535": {
			input:               types.StringValue("1-65535"),
			expectedDiagnostics: nil,
		},
		"valid-range-3000-4000": {
			input:               types.StringValue("3000-4000"),
			expectedDiagnostics: nil,
		},
		"valid-range-100-200": {
			input:               types.StringValue("100-200"),
			expectedDiagnostics: nil,
		},
		// Invalid port numbers
		"invalid-port-0": {
			input: types.StringValue("0"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be between 1 and 65535, got: 0",
				),
			},
		},
		"invalid-port-65536": {
			input: types.StringValue("65536"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be between 1 and 65535, got: 65536",
				),
			},
		},
		"invalid-port-99999": {
			input: types.StringValue("99999"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be between 1 and 65535, got: 99999",
				),
			},
		},
		"invalid-port-negative": {
			input: types.StringValue("-1"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: -1",
				),
			},
		},
		// Invalid formats
		"invalid-format-abc": {
			input: types.StringValue("abc"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be a valid number, got: abc",
				),
			},
		},
		"invalid-format-port80": {
			input: types.StringValue("port80"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be a valid number, got: port80",
				),
			},
		},
		"invalid-format-empty": {
			input: types.StringValue(""),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be a valid number, got: ",
				),
			},
		},
		"invalid-format-whitespace": {
			input: types.StringValue("  "),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Number",
					"Port must be a valid number, got:   ",
				),
			},
		},
		// Invalid ranges - start > end
		"invalid-range-start-greater-than-end": {
			input: types.StringValue("9000-8000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Start port must be less than end port, got: 9000-8000",
				),
			},
		},
		// Invalid ranges - start == end
		"invalid-range-start-equals-end": {
			input: types.StringValue("8000-8000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Start port must be less than end port, got: 8000-8000",
				),
			},
		},
		// Invalid range formats
		"invalid-range-format-no-end": {
			input: types.StringValue("8000-"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: 8000-",
				),
			},
		},
		"invalid-range-format-no-start": {
			input: types.StringValue("-9000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: -9000",
				),
			},
		},
		"invalid-range-format-double-hyphen": {
			input: types.StringValue("8000--9000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: 8000--9000",
				),
			},
		},
		"invalid-range-format-with-spaces": {
			input: types.StringValue("8000 - 9000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: 8000 - 9000",
				),
			},
		},
		// Invalid ranges - out of range start port
		"invalid-range-start-port-0": {
			input: types.StringValue("0-100"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Start port must be between 1 and 65535, got: 0",
				),
			},
		},
		"invalid-range-start-port-too-high": {
			input: types.StringValue("70000-80000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Start port must be between 1 and 65535, got: 70000",
				),
			},
		},
		// Invalid ranges - out of range end port
		"invalid-range-end-port-too-high": {
			input: types.StringValue("60000-70000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"End port must be between 1 and 65535, got: 70000",
				),
			},
		},
		"invalid-range-end-port-0": {
			input: types.StringValue("100-0"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"End port must be between 1 and 65535, got: 0",
				),
			},
		},
		// Invalid range with non-numeric values
		"invalid-range-non-numeric-start": {
			input: types.StringValue("abc-9000"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: abc-9000",
				),
			},
		},
		"invalid-range-non-numeric-end": {
			input: types.StringValue("8000-xyz"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Port Range",
					"Port range must be in format 'start-end' (e.g., '8000-9000'), got: 8000-xyz",
				),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    testCase.input,
			}
			response := validator.StringResponse{}
			IsPortOrRange().ValidateString(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
