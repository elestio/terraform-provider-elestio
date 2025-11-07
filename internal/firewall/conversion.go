package firewall

import (
	"context"
	"strings"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var toolRuleAttrTypes = map[string]attr.Type{
	"name":     types.StringType,
	"type":     types.StringType,
	"port":     types.StringType,
	"protocol": types.StringType,
	"targets":  types.SetType{ElemType: types.StringType},
}

// NormalizeTypeToTerraform converts API firewall type (uppercase) to Terraform format (lowercase)
func NormalizeTypeToTerraform(apiType string) string {
	return strings.ToLower(apiType)
}

// NormalizeTypeToAPI converts Terraform firewall type (lowercase) to API format (uppercase)
func NormalizeTypeToAPI(terraformType string) string {
	return strings.ToUpper(terraformType)
}

// ConvertElestioToTerraform converts Elestio API firewall rules to Terraform set
func ConvertElestioToTerraform(rules []elestio.ServiceFirewallRule, diags *diag.Diagnostics) types.Set {
	var firewallRulesObjs []attr.Value

	for _, rule := range rules {
		ruleObj := ConvertRuleToTerraformObject(rule, diags)
		if diags.HasError() {
			return types.SetNull(types.ObjectType{AttrTypes: models.FirewallRuleAttrTypes})
		}
		firewallRulesObjs = append(firewallRulesObjs, ruleObj)
	}

	firewallRulesSet, d := types.SetValue(types.ObjectType{AttrTypes: models.FirewallRuleAttrTypes}, firewallRulesObjs)
	diags.Append(d...)

	return firewallRulesSet
}

// ConvertRuleToTerraformObject converts a single Elestio rule to a Terraform object
func ConvertRuleToTerraformObject(rule elestio.ServiceFirewallRule, diags *diag.Diagnostics) types.Object {
	var targets []attr.Value
	for _, target := range rule.Targets {
		targets = append(targets, types.StringValue(target))
	}

	targetsSet, d := types.SetValue(types.StringType, targets)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(models.FirewallRuleAttrTypes)
	}

	ruleObj, d := types.ObjectValue(models.FirewallRuleAttrTypes, map[string]attr.Value{
		"type":     types.StringValue(NormalizeTypeToTerraform(rule.Type)),
		"port":     types.StringValue(rule.Port),
		"protocol": types.StringValue(rule.Protocol),
		"targets":  targetsSet,
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(models.FirewallRuleAttrTypes)
	}

	return ruleObj
}

// ConvertToolRuleToTerraformObject converts a single Elestio rule to a Terraform object with tool name
func ConvertToolRuleToTerraformObject(rule elestio.ServiceFirewallRule, name string, diags *diag.Diagnostics) types.Object {
	var targets []attr.Value
	for _, target := range rule.Targets {
		targets = append(targets, types.StringValue(target))
	}

	targetsSet, d := types.SetValue(types.StringType, targets)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(toolRuleAttrTypes)
	}

	ruleObj, d := types.ObjectValue(toolRuleAttrTypes, map[string]attr.Value{
		"name":     types.StringValue(name),
		"type":     types.StringValue(NormalizeTypeToTerraform(rule.Type)),
		"port":     types.StringValue(rule.Port),
		"protocol": types.StringValue(rule.Protocol),
		"targets":  targetsSet,
	})
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectNull(toolRuleAttrTypes)
	}

	return ruleObj
}

// ConvertTerraformToElestio converts Terraform firewall rules to Elestio API format
func ConvertTerraformToElestio(ctx context.Context, rulesSet types.Set, diags *diag.Diagnostics) []elestio.ServiceFirewallRule {
	var rules []models.FirewallRuleModel
	diags.Append(rulesSet.ElementsAs(ctx, &rules, true)...)
	if diags.HasError() {
		return nil
	}

	var elestioRules []elestio.ServiceFirewallRule
	for _, rule := range rules {
		var targets []string
		diags.Append(rule.Targets.ElementsAs(ctx, &targets, true)...)
		if diags.HasError() {
			return nil
		}

		elestioRules = append(elestioRules, elestio.ServiceFirewallRule{
			Type:     NormalizeTypeToAPI(rule.Type.ValueString()),
			Port:     rule.Port.ValueString(),
			Protocol: rule.Protocol.ValueString(),
			Targets:  targets,
		})
	}

	return elestioRules
}

// ExtractUserRules filters API rules, excluding API-managed tool ports to prevent state drift
func ExtractUserRules(ctx context.Context, apiRules []elestio.ServiceFirewallRule, planUserRules types.Set, diags *diag.Diagnostics) types.Set {
	planRulesMap := buildFirewallRulesKeySet(ctx, planUserRules, diags)

	var userRulesObjs []attr.Value
	for _, rule := range apiRules {
		port := rule.Port
		protocol := rule.Protocol
		ruleType := NormalizeTypeToTerraform(rule.Type)
		key := MakeFirewallRuleKey(port, protocol, ruleType)

		// Skip API-managed tool ports (not explicitly in user config)
		if IsToolPort(port, protocol, ruleType) && !planRulesMap[key] {
			continue
		}

		ruleObj := ConvertRuleToTerraformObject(rule, diags)
		if !diags.HasError() {
			userRulesObjs = append(userRulesObjs, ruleObj)
		}
	}

	userRulesSet, d := types.SetValue(types.ObjectType{AttrTypes: models.FirewallRuleAttrTypes}, userRulesObjs)
	diags.Append(d...)

	return userRulesSet
}

// ExtractToolRules returns API-managed tool ports not explicitly defined by user
func ExtractToolRules(ctx context.Context, apiRules []elestio.ServiceFirewallRule, planUserRules types.Set, diags *diag.Diagnostics) types.Set {
	planRulesMap := buildFirewallRulesKeySet(ctx, planUserRules, diags)

	var toolRulesObjs []attr.Value
	for _, rule := range apiRules {
		port := rule.Port
		protocol := rule.Protocol
		ruleType := NormalizeTypeToTerraform(rule.Type)
		key := MakeFirewallRuleKey(port, protocol, ruleType)

		// Only include tool ports that are NOT user-managed
		if IsToolPort(port, protocol, ruleType) && !planRulesMap[key] {
			toolName := GetToolName(port, protocol, ruleType)
			ruleObj := ConvertToolRuleToTerraformObject(rule, toolName, diags)
			if !diags.HasError() {
				toolRulesObjs = append(toolRulesObjs, ruleObj)
			}
		}
	}

	toolRulesSet, d := types.SetValue(types.ObjectType{AttrTypes: toolRuleAttrTypes}, toolRulesObjs)
	diags.Append(d...)

	return toolRulesSet
}

func buildFirewallRulesMap(ctx context.Context, rulesSet types.Set, diags *diag.Diagnostics) map[string]models.FirewallRuleModel {
	rulesMap := make(map[string]models.FirewallRuleModel)

	if rulesSet.IsNull() || rulesSet.IsUnknown() {
		return rulesMap
	}

	var rules []models.FirewallRuleModel
	diags.Append(rulesSet.ElementsAs(ctx, &rules, true)...)
	if diags.HasError() {
		return rulesMap
	}

	for _, rule := range rules {
		key := MakeFirewallRuleKey(
			rule.Port.ValueString(),
			rule.Protocol.ValueString(),
			rule.Type.ValueString(),
		)
		rulesMap[key] = rule
	}

	return rulesMap
}

func buildFirewallRulesKeySet(ctx context.Context, rulesSet types.Set, diags *diag.Diagnostics) map[string]bool {
	rulesMap := buildFirewallRulesMap(ctx, rulesSet, diags)
	keySet := make(map[string]bool, len(rulesMap))
	for key := range rulesMap {
		keySet[key] = true
	}
	return keySet
}
