package validators

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isPublicSSHKeyValidator struct{}

func (v isPublicSSHKeyValidator) Description(ctx context.Context) string {
	return "string should be a valid SSH key"
}

func (v isPublicSSHKeyValidator) MarkdownDescription(ctx context.Context) string {
	return "string should be a valid SSH key"
}

func (v isPublicSSHKeyValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	enteredKey := req.ConfigValue.ValueString()
	parts := strings.Split(enteredKey, " ")
	if len(parts) != 2 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			"The SSH public key should only contain two parts separated by a space."+
				" Example: `ssh-rsa AAaCfa...WAqDUNs=`."+
				" You should not include the username, hostname, or comment.",
		)
	}
}

func IsPublicSSHKey() isPublicSSHKeyValidator {
	return isPublicSSHKeyValidator{}
}
