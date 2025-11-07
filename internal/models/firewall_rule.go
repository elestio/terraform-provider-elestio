package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallRuleModel struct {
	Type     types.String `tfsdk:"type"`
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Targets  types.Set    `tfsdk:"targets"`
}

var FirewallRuleAttrTypes = map[string]attr.Type{
	"type":     types.StringType,
	"port":     types.StringType,
	"protocol": types.StringType,
	"targets":  types.SetType{ElemType: types.StringType},
}
