package firewall

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type firewallRulesConditionalValidator struct {
	defaultRules           types.Set
	defaultFirewallEnabled bool
}

func (v firewallRulesConditionalValidator) Description(ctx context.Context) string {
	return "When firewall_enabled is false, firewall_user_rules must be empty. When firewall_enabled is true, firewall_user_rules must contain required system input ports: 22/tcp/input (SSH) and 4242/udp/input (Nebula Global IP)"
}

func (v firewallRulesConditionalValidator) MarkdownDescription(ctx context.Context) string {
	return "When `firewall_enabled` is `false`, `firewall_user_rules` must be empty `[]`. When `firewall_enabled` is `true`, `firewall_user_rules` must contain required system input ports: 22/tcp/input (SSH) and 4242/udp/input (Nebula Global IP)"
}

func (v firewallRulesConditionalValidator) ValidateSet(ctx context.Context, req validator.SetRequest, resp *validator.SetResponse) {
	// Get the firewall_enabled attribute value
	var firewallEnabled types.Bool
	diags := req.Config.GetAttribute(ctx, path.Root("firewall_enabled"), &firewallEnabled)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If both firewall_enabled and firewall_user_rules are unknown/null, skip validation
	// The schema defaults will be applied after validation and will handle this correctly
	if (firewallEnabled.IsUnknown() || firewallEnabled.IsNull()) &&
		(req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull()) {
		return
	}

	// Determine if firewall is enabled (use configured value or default)
	isFirewallEnabled := v.defaultFirewallEnabled
	if !firewallEnabled.IsNull() && !firewallEnabled.IsUnknown() {
		isFirewallEnabled = firewallEnabled.ValueBool()
	}

	// Check if firewall_user_rules is empty
	isRulesEmpty := req.ConfigValue.IsNull() || len(req.ConfigValue.Elements()) == 0

	// Case 1: Firewall is disabled (explicitly or by default)
	if !isFirewallEnabled {
		// When firewall is explicitly disabled, firewall_user_rules MUST be explicitly set to empty
		// We cannot allow it to be unknown/null because that would result in default rules being applied
		if (!firewallEnabled.IsUnknown() && !firewallEnabled.IsNull()) &&
			(req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull()) {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Attribute Configuration",
				"When firewall_enabled is false, firewall_user_rules must be explicitly set to [].\n"+
					"Add: firewall_user_rules = []",
			)
			return
		}

		// If user specified rules, they must be empty
		if !isRulesEmpty {
			resp.Diagnostics.AddAttributeError(
				req.Path,
				"Invalid Attribute Configuration",
				"When firewall_enabled is false, firewall_user_rules must be empty [].\n"+
					"Remove all firewall rules or set firewall_enabled to true.",
			)
		}
		// No further validation needed when firewall is disabled
		return
	}

	// Case 2: Firewall is enabled (explicitly or by default)
	// Determine what value to validate: explicit config or default
	var valuesToValidate types.Set
	if req.ConfigValue.IsUnknown() || req.ConfigValue.IsNull() {
		// Use the default value provided by the schema
		valuesToValidate = v.defaultRules
	} else {
		// Use the explicitly provided value
		valuesToValidate = req.ConfigValue
	}

	// Validate that required system ports are present
	// Note: Template-specific ports are optional - only system ports are strictly required
	requiredPorts := GetSystemRequiredPorts()
	ValidateRequiredPorts(
		ctx,
		valuesToValidate,
		requiredPorts,
		&resp.Diagnostics,
		req.Path,
		"",
	)
}

func FirewallRulesConditional(defaultRules types.Set, defaultFirewallEnabled bool) firewallRulesConditionalValidator {
	return firewallRulesConditionalValidator{
		defaultRules:           defaultRules,
		defaultFirewallEnabled: defaultFirewallEnabled,
	}
}
