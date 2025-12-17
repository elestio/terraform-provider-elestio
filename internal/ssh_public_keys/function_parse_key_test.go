package ssh_public_keys

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func newParseKeyResultObject(keyData, username string) types.Object {
	result, _ := types.ObjectValue(
		map[string]attr.Type{
			"key_data": types.StringType,
			"username": types.StringType,
		},
		map[string]attr.Value{
			"key_data": types.StringValue(keyData),
			"username": types.StringValue(username),
		},
	)
	return result
}

func newParseKeyUnknownResult() types.Object {
	return types.ObjectUnknown(map[string]attr.Type{
		"key_data": types.StringType,
		"username": types.StringType,
	})
}

func TestParseKeyFunctionRun(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		request  function.RunRequest
		expected function.RunResponse
	}{
		"key-with-provided-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh"),
					types.StringValue("admin"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh",
					"admin",
				)),
			},
		},
		"key-with-comment-and-provided-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh user@host"),
					types.StringValue("admin"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh",
					"admin",
				)),
			},
		},
		"key-with-comment-null-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh user@host"),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh",
					"user@host",
				)),
			},
		},
		"key-with-comment-unknown-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh myuser"),
					types.StringUnknown(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh",
					"myuser",
				)),
			},
		},
		"ed25519-key-with-provided-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W"),
					types.StringValue("deploy"),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
					"deploy",
				)),
			},
		},
		"ed25519-key-with-comment-null-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W mykey@machine"),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
					"mykey@machine",
				)),
			},
		},
		"ecdsa-key-with-comment-null-username": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBA95ywHY2HQsFe59iIhJCNmPjQdGbAJ7/5ZcxfOdHs98gG6UhCj5KwjpSICNGTZ+ZE+W4ExRPWzAGfFzjibUzsE= cicd"),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBA95ywHY2HQsFe59iIhJCNmPjQdGbAJ7/5ZcxfOdHs98gG6UhCj5KwjpSICNGTZ+ZE+W4ExRPWzAGfFzjibUzsE=",
					"cicd",
				)),
			},
		},
		"key-with-leading-whitespace": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("  ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W test"),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
					"test",
				)),
			},
		},
		"key-with-trailing-whitespace": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W test   "),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
					"test",
				)),
			},
		},
		"key-with-trailing-newline": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W test\n"),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
					"test",
				)),
			},
		},
		"key-without-comment-null-username-error": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh"),
					types.StringNull(),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Username is required but not provided, and the SSH key does not contain a comment to use as username"),
				Result: function.NewResultData(newParseKeyUnknownResult()),
			},
		},
		"key-without-comment-unknown-username-error": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC0pA4vzGH+cmR+blZnoxO5HorOP1ubD4SxuOiW2DSNTSptlj+mPmFIL6sZeYMvSqAjXK368qL3DKHLpp2+1ws1XnYn/Zx/O4WBQAY7VbtzwFc7w7uirQaK6lVqXn8q4CnO0+5IYHgKLrNMEipwLKo+R3E3e1KrH5Xbyhj5yJzrMe3lWOAPzS27DJvjpN5SGWo65X6qFJRh3q95xOQhSOaEqZ/A2ZtfOuagq3FmASzoo/pbq7ianvnxzAYsb2Hg/9uAvypj4Beli6BP7419aP14XS0yyiW4XTKY/9XZiR/3VIKBN/stGN5NFLw82/j12E1GznbDG9PL7PQhijP7QgJh"),
					types.StringUnknown(),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Username is required but not provided, and the SSH key does not contain a comment to use as username"),
				Result: function.NewResultData(newParseKeyUnknownResult()),
			},
		},
		"invalid-key-only-type": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-rsa"),
					types.StringValue("admin"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)"),
				Result: function.NewResultData(newParseKeyUnknownResult()),
			},
		},
		"invalid-key-empty": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue(""),
					types.StringValue("admin"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)"),
				Result: function.NewResultData(newParseKeyUnknownResult()),
			},
		},
		"invalid-key-whitespace": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("   "),
					types.StringValue("admin"),
				}),
			},
			expected: function.RunResponse{
				Error:  function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)"),
				Result: function.NewResultData(newParseKeyUnknownResult()),
			},
		},
		"empty-username-provided": {
			request: function.RunRequest{
				Arguments: function.NewArgumentsData([]attr.Value{
					types.StringValue("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W"),
					types.StringValue(""),
				}),
			},
			expected: function.RunResponse{
				Result: function.NewResultData(newParseKeyResultObject(
					"ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOwlR9xtbM69hWLJbB5nHi0a65TuRvtaldgTJQ4ClL1W",
					"",
				)),
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := function.RunResponse{
				Result: function.NewResultData(newParseKeyUnknownResult()),
			}

			(&ParseKeyFunction{}).Run(context.Background(), testCase.request, &got)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
