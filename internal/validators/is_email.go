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

	_, err := mail.ParseAddress(value)
	if err != nil {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Email Address",
			fmt.Sprintf("Invalid email address: %s", value),
		)
	}
}

func IsEmail() isEmailValidator {
	return isEmailValidator{}
}
