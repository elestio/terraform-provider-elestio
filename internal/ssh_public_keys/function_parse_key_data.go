package ssh_public_keys

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = &ParseKeyDataFunction{}

type ParseKeyDataFunction struct{}

func NewParseKeyDataFunction() function.Function {
	return &ParseKeyDataFunction{}
}

func (f *ParseKeyDataFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse_ssh_key_data"
}

func (f *ParseKeyDataFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Strips the comment from an SSH public key",
		Description: "Takes an SSH public key string and returns only the key type and key data, removing any trailing comment. This is useful because Elestio does not accept SSH keys with comments.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "ssh_key",
				Description: "The SSH public key string (e.g., 'ssh-rsa AAAA... user@host')",
			},
		},
		Return: function.StringReturn{},
	}
}

func (f *ParseKeyDataFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sshKey string

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sshKey))
	if resp.Error != nil {
		return
	}

	// Trim whitespace and newlines
	sshKey = strings.TrimSpace(sshKey)

	fields := strings.Fields(sshKey)
	if len(fields) < 2 {
		resp.Error = function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)")
		return
	}

	// Return only the first two fields (key type and key data)
	result := fields[0] + " " + fields[1]
	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}

