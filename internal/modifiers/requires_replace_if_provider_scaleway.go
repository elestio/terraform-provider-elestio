package modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func RequiresReplaceIfProviderScaleway() stringplanmodifier.RequiresReplaceIfFunc {
	return func(ctx context.Context, req planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
		if req.ConfigValue.IsUnknown() || req.StateValue.IsNull() {
			resp.RequiresReplace = false
			return
		}

		var providerName string
		req.Config.GetAttribute(ctx, path.Root("provider_name"), &providerName)
		if providerName != "scaleway" {
			resp.RequiresReplace = false
			return
		}

		resp.RequiresReplace = true
	}
}
