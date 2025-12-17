package ssh_public_keys

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsValidKey(t *testing.T) {
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
		"empty": {
			input: types.StringValue(""),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Expected a non-empty or whitespace string.",
				),
			},
		},
		"whitespace": {
			input: types.StringValue("   "),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Expected a non-empty or whitespace string.",
				),
			},
		},
		"newline-n": {
			input: types.StringValue("ssh \n newline"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Your SSH public key must be on a single line. "+
						"Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to read and parse your key file. "+
						"Read the guide: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys",
				),
			},
		},
		"newline-r": {
			input: types.StringValue("ssh \r newline"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Your SSH public key must be on a single line. "+
						"Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to read and parse your key file. "+
						"Read the guide: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys",
				),
			},
		},
		"valid-ecdsa-sha2-nistp256": {
			input:               types.StringValue("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBA95ywHY2HQsFe59iIhJCNmPjQdGbAJ7/5ZcxfOdHs98gG6UhCj5KwjpSICNGTZ+ZE+W4ExRPWzAGfFzjibUzsE="),
			expectedDiagnostics: nil,
		},
		"valid-ssh-ed25519": {
			input:               types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W"),
			expectedDiagnostics: nil,
		},
		"incomplete-key": {
			input: types.StringValue("ssh-rsa"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Invalid SSH key format: expected at least key type and key data separated by a space. "+
						"Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to read and parse your key file. "+
						"Read the guide: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys",
				),
			},
		},
		"rsa-fake-key": {
			input: types.StringValue("ssh-rsa ThisIsNot a REAL key"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Error decoding the key data. Your SSH public key does not seem to be base64 encoded. "+
						"Read the guide: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys",
				),
			},
		},
		"rsa-altered-key": {
			input: types.StringValue("ssh-rsa AAAAB3NzbC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"Error parsing the key data. Your SSH public key does not seem to be valid. It may be corrupted or altered. "+
						"Read the guide: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys",
				),
			},
		},
		"good-key": {
			// 2048 bits RSA key
			input:               types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh"),
			expectedDiagnostics: nil,
		},
		"good-key-with-comment": {
			// 2048 bits RSA key with comment
			input: types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh user@host"),
			expectedDiagnostics: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid Attribute Configuration",
					"SSH public key comments are not allowed. "+
						"Use `provider::elestio::parse_ssh_key_data()` to remove the comment, "+
						"or `provider::elestio::parse_ssh_key()` to extract both the key and username from the comment. "+
						"Read the guide: https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys",
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
			IsValidKey().ValidateString(context.Background(), request, &response)

			if diff := cmp.Diff(response.Diagnostics, testCase.expectedDiagnostics); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}
