package modifiers

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RemoveSSHKeyCommentModifier struct{}

func (m RemoveSSHKeyCommentModifier) Description(ctx context.Context) string {
	return "Remove the comment from the ssh key if it exist"
}

func (m RemoveSSHKeyCommentModifier) MarkdownDescription(ctx context.Context) string {
	return "Remove the comment from the ssh key if it exist"
}

func (m RemoveSSHKeyCommentModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		key := req.PlanValue.ValueString()
		keyParts := strings.Split(key, " ")
		if len(keyParts) > 2 {
			key = strings.Join(keyParts[:2], " ")
			resp.PlanValue = types.StringValue(key)
		}
	}

}

func RemoveSSHKeyComment() RemoveSSHKeyCommentModifier {
	return RemoveSSHKeyCommentModifier{}
}
