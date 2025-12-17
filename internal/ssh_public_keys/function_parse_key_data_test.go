package ssh_public_keys

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestParseKeyDataFunctionRun(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"valid-key-without-comment": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh")),
			},
		},
		"valid-key-with-comment": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh user@host"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh")),
			},
		},
		"valid-ed25519-key": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W")),
			},
		},
		"valid-ed25519-key-with-comment": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W mykey@machine"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W")),
			},
		},
		"valid-ecdsa-key-with-comment": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBA95ywHY2HQsFe59iIhJCNmPjQdGbAJ7/5ZcxfOdHs98gG6UhCj5KwjpSICNGTZ+ZE+W4ExRPWzAGfFzjibUzsE= admin"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBA95ywHY2HQsFe59iIhJCNmPjQdGbAJ7/5ZcxfOdHs98gG6UhCj5KwjpSICNGTZ+ZE+W4ExRPWzAGfFzjibUzsE=")),
			},
		},
		"key-with-leading-whitespace": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("  ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W")),
			},
		},
		"key-with-trailing-whitespace": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W   "),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W")),
			},
		},
		"key-with-trailing-newline": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W\n"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W")),
			},
		},
		"invalid-key-only-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"invalid-key-empty": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(""),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
		"invalid-key-whitespace": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("   "),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)"),
				Result: function.NewResultData(types.StringUnknown()),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := function.RunResponse{
				Result: function.NewResultData(types.StringUnknown()),
			}

			(&ParseKeyDataFunction{}).Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
