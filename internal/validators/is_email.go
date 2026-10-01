package validators

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isEmailValidator struct{}

func (v isEmailValidator) Description(ctx context.Context) string {
	return "string should be a valid email address"
}

func (v isEmailValidator) MarkdownDescription(ctx context.Context) string {
	return "string should be a valid email address"
}

func (v isEmailValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	// Only a bare address is accepted. mail.ParseAddress alone also accepts
	// display-name forms such as "Name <a@b.c>" and comments.
	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Address != value {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Email Address",
			fmt.Sprintf("Invalid email address: %q", value),
		)
	}
}

func IsEmail() isEmailValidator {
	return isEmailValidator{}
}
