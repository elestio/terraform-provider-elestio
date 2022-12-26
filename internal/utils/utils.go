package utils

import (
	"fmt"

	"github.com/elestio/elestio-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MapStringToMapType maps a native golang map to a terraform map type
func MapStringToMapType(m map[string]string, diags *diag.Diagnostics) types.Map {
	elems := make(map[string]attr.Value)
	for k, v := range m {
		elems[k] = types.StringValue(v)
	}
	v, d := types.MapValue(types.StringType, elems)
	diags.Append(d...)
	return v
}

// ObjectValue is a helper function to build a terraform object type
func ObjectValue(attrTypes map[string]attr.Type, attrs map[string]attr.Value, diags *diag.Diagnostics) types.Object {
	object, err := types.ObjectValue(attrTypes, attrs)
	if err != nil {
		diags.AddError(
			"Error Building Object Type",
			fmt.Sprintf("Unable to build the object type, got error: %s", err),
		)
	}
	return object
}

// BoolValue is a helper function to build a terraform bool type
func BoolValue(elestioBool elestio.NumberAsBool) types.Bool {
	if elestioBool == elestio.NumberAsBool(1) {
		return types.BoolValue(true)
	}

	return types.BoolValue(false)
}
