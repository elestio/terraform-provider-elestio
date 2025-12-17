package provider

import (
	ssh_public_keys "github.com/elestio/terraform-provider-elestio/internal/ssh_public_keys"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ====================================================================================
// DEPRECATED CODE - TO BE REMOVED IN NEXT MAJOR VERSION
// ====================================================================================
// This file contains all deprecated code that should be removed in the next major update.

// SSHKeyModel - Deprecated: replaced by models.SSHPublicKeyModel
type SSHKeyModel struct {
	KeyName   types.String `tfsdk:"key_name"`
	PublicKey types.String `tfsdk:"public_key"`
}

// sshKeyAttrTypes - Deprecated: replaced by models.SSHPublicKeyAttrTypes
var sshKeyAttrTypes = map[string]attr.Type{
	"key_name":   types.StringType,
	"public_key": types.StringType,
}

// SSHKeysDefaultValue - Deprecated: replaced by SSHPublicKeysDefaultValue() from ssh_public_keys.go
func SSHKeysDefaultValue() types.Set {
	set, _ := types.SetValue(types.ObjectType{AttrTypes: sshKeyAttrTypes}, []attr.Value{})
	return set
}

// sshKeysSchema - Deprecated: replaced by sshPublicKeysSchema from ssh_public_keys.go
var sshKeysSchema = schema.SetNestedAttribute{
	MarkdownDescription: "This attribute allows you to add SSH keys to your service.",
	DeprecationMessage:  "This attribute is deprecated and will be removed in a future version. Please use `ssh_public_keys` instead.",
	Optional:            true,
	Computed:            true,
	Default:             setdefault.StaticValue(SSHKeysDefaultValue()),
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"key_name": schema.StringAttribute{
				MarkdownDescription: "SSH Key Name.",
				Required:            true,
			},
			"public_key": schema.StringAttribute{
				MarkdownDescription: "SSH Public Key." +
					"The SSH public key should only contain two parts separated by a space." +
					" Example: `ssh-rsa AAaCfa...WAqDUNs=`." +
					" You should not include the username, hostname, or comment.",
				Required: true,
				Validators: []validator.String{
					ssh_public_keys.IsValidKey(),
				},
				DeprecationMessage: "This attribute is deprecated and will be removed in a future version." +
					" Please use `ssh_public_keys` instead.",
			},
		},
	},
}
