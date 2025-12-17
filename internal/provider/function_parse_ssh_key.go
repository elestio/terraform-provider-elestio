package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = &ParseSSHKeyFunction{}

type ParseSSHKeyFunction struct{}

func NewParseSSHKeyFunction() function.Function {
	return &ParseSSHKeyFunction{}
}

func (f *ParseSSHKeyFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "parse_ssh_key"
}

func (f *ParseSSHKeyFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Parses an SSH public key and extracts key data and username",
		Description: "Takes an SSH public key and a username parameter. Returns an object with `key_data` (the key without comment) and `username` (from the provided value or extracted from the key comment). The username parameter is required but accepts null - pass null explicitly to extract the username from the SSH key comment. If username is null and the key has no comment, an error is returned.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "ssh_key",
				Description: "The SSH public key string (e.g., 'ssh-rsa AAAA... user@host')",
			},
			function.StringParameter{
				Name:               "username",
				Description:        "The username to use. Pass null explicitly to extract the username from the SSH key comment.",
				AllowNullValue:     true,
				AllowUnknownValues: true,
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: map[string]attr.Type{
				"key_data": types.StringType,
				"username": types.StringType,
			},
		},
	}
}

func (f *ParseSSHKeyFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var sshKey types.String
	var username types.String

	resp.Error = function.ConcatFuncErrors(resp.Error, req.Arguments.Get(ctx, &sshKey, &username))
	if resp.Error != nil {
		return
	}

	// Trim whitespace and newlines
	sshKeyStr := strings.TrimSpace(sshKey.ValueString())

	fields := strings.Fields(sshKeyStr)
	if len(fields) < 2 {
		resp.Error = function.NewFuncError("Invalid SSH key format: expected at least 2 space-separated fields (key type and key data)")
		return
	}

	// Build key_data (first two fields only)
	keyData := fields[0] + " " + fields[1]

	// Determine username
	var resultUsername string
	if !username.IsNull() && !username.IsUnknown() {
		// Username is provided, use it
		resultUsername = username.ValueString()
	} else {
		// Username is null or unknown, try to extract from key comment
		if len(fields) >= 3 {
			resultUsername = fields[2]
		} else {
			resp.Error = function.NewFuncError("Username is required but not provided, and the SSH key does not contain a comment to use as username")
			return
		}
	}

	// Build result object
	result, diags := types.ObjectValue(
		map[string]attr.Type{
			"key_data": types.StringType,
			"username": types.StringType,
		},
		map[string]attr.Value{
			"key_data": types.StringValue(keyData),
			"username": types.StringValue(resultUsername),
		},
	)

	if diags.HasError() {
		resp.Error = function.NewFuncError("Failed to create result object")
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
