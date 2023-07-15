package validators

import (
	"context"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isDefaultPasswordValidator struct{}

func (v isDefaultPasswordValidator) Description(ctx context.Context) string {
	return "default password should respect the format"
}

func (v isDefaultPasswordValidator) MarkdownDescription(ctx context.Context) string {
	return "default password should respect the format"
}

func (v isDefaultPasswordValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	enteredPassword := req.ConfigValue.ValueString()

	// Password length must be at least 10 characters
	if len(enteredPassword) < 10 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			"The password must be at least 10 characters.",
		)
	}

	// Password can only contain alphanumeric characters or hyphens
	if !regexp.MustCompile(`^[a-zA-Z0-9-]+$`).MatchString(enteredPassword) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			"The password can only contain alphanumeric characters or hyphens.",
		)
	}

	// Password must contain at least one uppercase letter
	if !strings.ContainsAny(enteredPassword, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			"The password must contain at least one uppercase letter.",
		)
	}

	// Password must contain at least one lowercase letter
	if !strings.ContainsAny(enteredPassword, "abcdefghijklmnopqrstuvwxyz") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			"The password must contain at least one lowercase letter.",
		)
	}

	// Password must contain at least one number
	if !strings.ContainsAny(enteredPassword, "0123456789") {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Attribute Configuration",
			"The password must contain at least one number.",
		)
	}
}

func IsDefaultPassword() isDefaultPasswordValidator {
	return isDefaultPasswordValidator{}
}
