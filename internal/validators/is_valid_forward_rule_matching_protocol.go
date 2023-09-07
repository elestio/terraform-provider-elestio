package validators

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/exp/slices"
)

var _ validator.String = &isValidForwardRuleMatchingProtocol{}

type isValidForwardRuleMatchingProtocol struct {
	expression path.Expression
	values     []types.String
}

func (v isValidForwardRuleMatchingProtocol) Description(_ context.Context) string {
	return fmt.Sprintf("If you want to use one of theses forward rules protocols %v, `protocol` and `target_protocol` must have the same value.", v.values)
}

func (v isValidForwardRuleMatchingProtocol) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v isValidForwardRuleMatchingProtocol) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	expression := req.PathExpression.Merge(v.expression)

	matchedPaths, diags := req.Config.PathMatches(ctx, expression)
	resp.Diagnostics.Append(diags...)
	// Still continue if there are errors, as we want to collect all of them.

	for _, matchedPath := range matchedPaths {
		var matchedPathValue attr.Value
		diags := req.Config.GetAttribute(ctx, matchedPath, &matchedPathValue)
		resp.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		// If the matched path value is null or unknown, we cannot compare
		// values, so continue to other matched paths.
		if matchedPathValue.IsNull() || matchedPathValue.IsUnknown() {
			continue
		}

		// Now that we know the matched path value is not null or unknown,
		// it is safe to attempt converting it to the intended attr.Value
		// implementation, in this case a types.String value.
		var matchedPathConfig types.String
		diags = tfsdk.ValueAs(ctx, matchedPathValue, &matchedPathConfig)
		resp.Diagnostics.Append(diags...)

		// If the matched path value was not able to be converted from
		// attr.Value to the intended types.String implementation, it most
		// likely means that the path expression was not pointing at a
		// types.StringType attribute. Collect the error and continue to
		// other matched paths.
		if diags.HasError() {
			continue
		}

		isSelfValueOneOf := slices.ContainsFunc(v.values, func(v types.String) bool {
			return v.Equal(req.ConfigValue)
		})

		if isSelfValueOneOf {
			isIdenticalToMe := req.ConfigValue.Equal(matchedPathConfig)
			if !isIdenticalToMe {
				resp.Diagnostics.AddAttributeError(
					req.Path,
					"Invalid Attribute Value",
					fmt.Sprintf("If you want to use one of theses forward rules protocols %v, `protocol` and `target_protocol` must have the same value.", v.values),
				)
			}
		}

	}
}

func IsValidForwardRuleMatchingProtocol(expression path.Expression, values ...string) validator.String {
	frameworkValues := make([]types.String, 0, len(values))

	for _, value := range values {
		frameworkValues = append(frameworkValues, types.StringValue(value))
	}

	return &isValidForwardRuleMatchingProtocol{
		expression: expression,
		values:     frameworkValues,
	}
}
