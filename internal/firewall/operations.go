package firewall

import (
	"context"
	"fmt"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Enable enables the firewall with user rules
func Enable(ctx context.Context, serviceID string, planRules types.Set, diags *diag.Diagnostics, client *elestio.Client) error {
	userRules := ConvertTerraformToElestio(ctx, planRules, diags)
	if diags.HasError() {
		return fmt.Errorf("failed to convert firewall rules")
	}

	if err := client.Service.EnableFirewallWithRules(serviceID, userRules); err != nil {
		return fmt.Errorf("failed to enable firewall with rules: %s", err)
	}

	return nil
}

// Update updates firewall rules, optionally removing API-managed tool ports
func Update(ctx context.Context, serviceID string, apiRules []elestio.ServiceFirewallRule, planRules types.Set, removeToolPorts bool, diags *diag.Diagnostics, client *elestio.Client) error {
	userRules := ConvertTerraformToElestio(ctx, planRules, diags)
	if diags.HasError() {
		return fmt.Errorf("failed to convert firewall rules")
	}

	mergedRules := MergeWithToolPorts(userRules, apiRules, removeToolPorts)

	if err := client.Service.UpdateFirewallRules(serviceID, mergedRules); err != nil {
		return fmt.Errorf("failed to update firewall rules: %s", err)
	}

	return nil
}

// ApplyChangesIfNeeded applies firewall changes only if rules or removeToolPorts changed
func ApplyChangesIfNeeded(
	ctx context.Context,
	serviceID string,
	apiRules []elestio.ServiceFirewallRule,
	stateUserRules types.Set,
	stateRemoveToolPorts types.Bool,
	planUserRules types.Set,
	planRemoveToolPorts types.Bool,
	diags *diag.Diagnostics,
	client *elestio.Client,
) error {
	rulesChanged := !stateUserRules.Equal(planUserRules)
	removeToolPortsChanged := !stateRemoveToolPorts.Equal(planRemoveToolPorts)

	needsUpdate := false

	if rulesChanged {
		needsUpdate = true
	} else if removeToolPortsChanged && planRemoveToolPorts.ValueBool() {
		// Only update if there are actual API-managed tool ports to remove
		needsUpdate = hasAPIManagedToolPorts(ctx, apiRules, planUserRules, diags)
		if diags.HasError() {
			return fmt.Errorf("failed to check for API-managed tool ports")
		}
	}

	if needsUpdate {
		return Update(ctx, serviceID, apiRules, planUserRules, planRemoveToolPorts.ValueBool(), diags, client)
	}

	return nil
}

// MergeWithToolPorts combines user rules with API tool ports (unless removeToolPorts is true)
func MergeWithToolPorts(userRules []elestio.ServiceFirewallRule, apiRules []elestio.ServiceFirewallRule, removeToolPorts bool) []elestio.ServiceFirewallRule {
	userRulesMap := make(map[string]bool)
	for _, rule := range userRules {
		key := MakeFirewallRuleKey(rule.Port, rule.Protocol, rule.Type)
		userRulesMap[key] = true
	}

	merged := append([]elestio.ServiceFirewallRule{}, userRules...)

	if !removeToolPorts {
		for _, apiRule := range apiRules {
			port := apiRule.Port
			protocol := apiRule.Protocol
			ruleType := NormalizeTypeToTerraform(apiRule.Type)
			key := MakeFirewallRuleKey(port, protocol, ruleType)

			if userRulesMap[key] || !IsToolPort(port, protocol, ruleType) {
				continue
			}

			merged = append(merged, apiRule)
		}
	}

	return merged
}

func hasAPIManagedToolPorts(ctx context.Context, apiRules []elestio.ServiceFirewallRule, planUserRules types.Set, diags *diag.Diagnostics) bool {
	planRulesMap := buildFirewallRulesKeySet(ctx, planUserRules, diags)
	if diags.HasError() {
		return false
	}

	for _, rule := range apiRules {
		port := rule.Port
		protocol := rule.Protocol
		ruleType := NormalizeTypeToTerraform(rule.Type)
		key := MakeFirewallRuleKey(port, protocol, ruleType)

		if IsToolPort(port, protocol, ruleType) && !planRulesMap[key] {
			return true
		}
	}

	return false
}
