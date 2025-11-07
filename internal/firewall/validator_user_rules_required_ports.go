package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type firewallUserRulesContainsRequiredPortsValidator struct{}

func (v firewallUserRulesContainsRequiredPortsValidator) Description(ctx context.Context) string {
	return "firewall_user_rules must contain required system input ports: 22/tcp/input (SSH) and 4242/udp/input (Nebula Global IP)"
}

func (v firewallUserRulesContainsRequiredPortsValidator) MarkdownDescription(ctx context.Context) string {
	return "firewall_user_rules must contain required system input ports: 22/tcp/input (SSH) and 4242/udp/input (Nebula Global IP)"
}

func (v firewallUserRulesContainsRequiredPortsValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// If the value is unknown or null, there is nothing to validate.
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		return
	}

	requiredPorts := GetSystemRequiredPorts()

	ValidateRequiredPorts(
		ctx,
		req.ConfigValue,
		requiredPorts,
		&resp.Diagnostics,
		req.Path,
		"",
	)
}

func FirewallUserRulesContainsRequiredPorts() firewallUserRulesContainsRequiredPortsValidator {
	return firewallUserRulesContainsRequiredPortsValidator{}
}
