package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SSHPublicKeyModel struct {
	Username types.String `tfsdk:"username"`
	KeyData  types.String `tfsdk:"key_data"`
}

var SSHPublicKeyAttrTypes = map[string]attr.Type{
	"username": types.StringType,
	"key_data": types.StringType,
}
