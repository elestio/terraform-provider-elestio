package ssh_public_keys

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/elestio/terraform-provider-elestio/internal/models"
)

// DefaultValue returns the default empty set for SSH public keys
func DefaultValue() types.Set {
	return types.SetValueMust(types.ObjectType{AttrTypes: models.SSHPublicKeyAttrTypes}, []attr.Value{})
}

// Schema returns the schema for ssh_public_keys attribute
var Schema = schema.SetNestedAttribute{
	MarkdownDescription: "You can add Public Keys to your resource to access it via the SSH protocol.",
	Optional:            true,
	Computed:            true,
	Default:             setdefault.StaticValue(DefaultValue()),
	Validators: []validator.Set{
		UniqueUsernames(),
	},
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "The username is used to identify the Public Key among others. Must be unique (per resource).",
				Required:            true,
			},
			"key_data": schema.StringAttribute{
				MarkdownDescription: "The Public Key value without comment." +
					" Use `provider::elestio::parse_ssh_key_data(file(\"~/.ssh/id_rsa.pub\"))` to remove the comment from your key." +
					" Read the guide [\"How to use SSH keys with Elestio Terraform Provider\"](https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/ssh_keys).",
				Required: true,
				Validators: []validator.String{
					IsValidKey(),
				},
			},
		},
	},
}
