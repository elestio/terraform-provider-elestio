package provider

import (
	"context"
	"fmt"

	"github.com/elestio/terraform-provider-elestio/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type SSHPublicKeyModel struct {
	Username types.String `tfsdk:"username"`
	KeyData  types.String `tfsdk:"key_data"`
}

var sshPublicKeyAttryTypes = map[string]attr.Type{
	"username": types.StringType,
	"key_data": types.StringType,
}

var sshPublicKeysDefaultValue = types.SetValueMust(types.ObjectType{AttrTypes: sshPublicKeyAttryTypes}, []attr.Value{})

var sshPublicKeysSchema = schema.SetNestedAttribute{
	MarkdownDescription: "You can add Public Keys to your resource to access it via the SSH protocol.",
	Optional:            true,
	Computed:            true,
	Default:             setdefault.StaticValue(sshPublicKeysDefaultValue),
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username is used to identify the Public Key among others. Must be unique (per resource).",
				Required:            true,
			},
			"key_data": schema.StringAttribute{
				MarkdownDescription: "The Public Key value." +
					" Read the guide [\"How generate a valid SSH Key for Elestio\"](https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys)." +
					" Example: `key_data = chomp(file('~/.ssh/id_rsa.pub'))`.",
				Required: true,
				Validators: []validator.String{
					validators.IsSSHPublicKey(),
				},
			},
		},
	},
}

// TODO: Move this check to a validator file with the proper format
func ensureSSHPublicKeysUsernamesAreUnique(ctx *context.Context, data *basetypes.SetValue, diags *diag.Diagnostics) {
	var keys []SSHPublicKeyModel

	diags.Append(data.ElementsAs(*ctx, &keys, true)...)
	if diags.HasError() {
		return
	}

	keysMap := make(map[string]SSHPublicKeyModel)
	for _, key := range keys {
		if duplicatedUsername, exists := keysMap[key.Username.ValueString()]; exists {
			diags.AddAttributeError(
				path.Root("ssh_public_keys"),
				"Invalid Attribute Configuration",
				"SSH Public Key Username must be unique per ressource."+
					fmt.Sprintf(" The following username is duplicated: %s", duplicatedUsername.Username.ValueString()),
			)
			return
		}
		keysMap[key.Username.ValueString()] = key
	}
}

func compareSSHPublicKeys(ctx *context.Context, state *basetypes.SetValue, plan *basetypes.SetValue, diags *diag.Diagnostics) (toAdd []SSHPublicKeyModel, toUpdate []SSHPublicKeyModel, toRemove []SSHPublicKeyModel) {
	var stateKeys []SSHPublicKeyModel
	diags.Append(state.ElementsAs(*ctx, &stateKeys, true)...)
	if diags.HasError() {
		return nil, nil, nil
	}
	var planKeys []SSHPublicKeyModel
	diags.Append(plan.ElementsAs(*ctx, &planKeys, true)...)
	if diags.HasError() {
		return nil, nil, nil
	}

	toAdd, toUpdate, toRemove = []SSHPublicKeyModel{}, []SSHPublicKeyModel{}, []SSHPublicKeyModel{}

	if len(stateKeys) == 0 && len(planKeys) == 0 {
		return toAdd, toUpdate, toRemove
	}

	// Usernames are unique and can be used as map index
	stateKeysMap := make(map[string]SSHPublicKeyModel)
	for _, obj := range stateKeys {
		stateKeysMap[obj.Username.ValueString()] = obj
	}
	planKeysMap := make(map[string]SSHPublicKeyModel)
	for _, obj := range planKeys {
		planKeysMap[obj.Username.ValueString()] = obj
	}

	// Iterate over state and delete any objects that are not in the plan
	for _, stateKey := range stateKeys {
		if _, exists := planKeysMap[stateKey.Username.ValueString()]; !exists {
			toRemove = append(toRemove, stateKey)
		}
	}

	// Iterate over the plan and compare each key to the corresponding key in state
	for _, planKey := range planKeys {
		if stateKey, exists := stateKeysMap[planKey.Username.ValueString()]; exists {
			if !planKey.KeyData.Equal(stateKey.KeyData) {
				toUpdate = append(toUpdate, planKey)
			}
		} else {
			toAdd = append(toAdd, planKey)
		}
	}

	return toAdd, toUpdate, toRemove
}
