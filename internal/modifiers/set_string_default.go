package modifiers

import (
	"context"
	"fmt"

	"github.com/elestio/terraform-provider-elestio/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type setStringDefaultModifier struct {
	Default *[]string
}

func (m setStringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %v", m.Default)
}

func (m setStringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %v", m.Default)
}

func (m setStringDefaultModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	if req.PlanValue.IsUnknown() || req.PlanValue.IsNull() {
		resp.PlanValue = utils.SliceStringToSetType(*m.Default, &resp.Diagnostics)
		return
	}
}

func SetStringDefault(def []string) setStringDefaultModifier {
	return setStringDefaultModifier{Default: &def}
}
