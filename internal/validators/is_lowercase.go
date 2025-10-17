package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type isLowercaseValidator struct{}

func (v isLowercaseValidator) Description(ctx context.Context) string {
	return "string must be lowercase"
}

func (v isLowercaseValidator) MarkdownDescription(ctx context.Context) string {
	return "string must be lowercase"
}

func (v isLowercaseValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	value := req.ConfigValue.ValueString()

	if value != strings.ToLower(value) {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Value",
			fmt.Sprintf("The value must be lowercase. Got '%s', expected '%s'.", value, strings.ToLower(value)),
		)
	}
}

func IsLowercase() isLowercaseValidator {
	return isLowercaseValidator{}
}
