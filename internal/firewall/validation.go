package firewall

import (
	"context"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ValidatePort80Required checks if port 80/tcp/input is present when custom domains are used
func ValidatePort80Required(ctx context.Context, firewallUserRules types.Set, diags *diag.Diagnostics, attributePath path.Path) bool {
	requiredPorts := GetCustomDomainRequiredPorts()

	return ValidateRequiredPorts(
		ctx,
		firewallUserRules,
		requiredPorts,
		diags,
		attributePath,
		"When custom domain names are specified,",
	)
}

// ValidateRulesForCustomDomains validates firewall rules when custom domains are present
// If firewallUserRules is null/unknown, validates against the default template rules
func ValidateRulesForCustomDomains(ctx context.Context, firewallUserRules types.Set, customDomainNames types.Set, firewallEnabled bool, templateRules []elestio.ServiceFirewallRule, hasCustomFirewallPorts bool, diags *diag.Diagnostics, attributePath path.Path) bool {
	if !firewallEnabled || customDomainNames.IsNull() || customDomainNames.IsUnknown() {
		return true
	}

	var customDomains []string
	d := customDomainNames.ElementsAs(ctx, &customDomains, false)
	diags.Append(d...)
	if diags.HasError() {
		return false
	}

	if len(customDomains) == 0 {
		return true
	}

	// If firewall_user_rules is null/unknown, use the default value for validation
	rulesToValidate := firewallUserRules
	if firewallUserRules.IsNull() || firewallUserRules.IsUnknown() {
		rulesToValidate = CreateUserRulesDefaultValue(templateRules, hasCustomFirewallPorts)
	}

	return ValidatePort80Required(ctx, rulesToValidate, diags, attributePath)
}
