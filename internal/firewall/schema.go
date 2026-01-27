package firewall

import (
	"fmt"
	"strings"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/models"
	"github.com/elestio/terraform-provider-elestio/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TargetsDefaultValue returns the default targets (all IPv4/IPv6)
func TargetsDefaultValue() types.Set {
	defaultTargets := GetDefaultTargets()
	targetsSet, _ := types.SetValue(types.StringType, []attr.Value{
		types.StringValue(defaultTargets[0]),
		types.StringValue(defaultTargets[1]),
	})
	return targetsSet
}

// CreateUserRulesDefaultValue creates the default firewall rules from template rules
func CreateUserRulesDefaultValue(templateRules []elestio.ServiceFirewallRule, hasCustomFirewallPorts bool) types.Set {
	// If firewall is disabled by default (no custom firewall ports), return empty array
	if !hasCustomFirewallPorts {
		emptySet, _ := types.SetValue(types.ObjectType{AttrTypes: models.FirewallRuleAttrTypes}, []attr.Value{})
		return emptySet
	}

	// If firewall is enabled by default, include required system ports + template rules
	var firewallRulesObjs []attr.Value
	diags := &diag.Diagnostics{}

	for _, requiredRule := range GetRequiredSystemPorts() {
		ruleObj := ConvertRuleToTerraformObject(requiredRule, diags)
		if !diags.HasError() {
			firewallRulesObjs = append(firewallRulesObjs, ruleObj)
		}
	}

	for _, rule := range templateRules {
		ruleObj := ConvertRuleToTerraformObject(rule, diags)
		if !diags.HasError() {
			firewallRulesObjs = append(firewallRulesObjs, ruleObj)
		}
	}

	firewallRulesSet, d := types.SetValue(types.ObjectType{AttrTypes: models.FirewallRuleAttrTypes}, firewallRulesObjs)
	diags.Append(d...)
	return firewallRulesSet
}

// UserRulesSchema returns the schema for firewall_user_rules
func UserRulesSchema(templateRules []elestio.ServiceFirewallRule, hasCustomFirewallPorts bool, supportsCustomDomains bool) schema.SetNestedAttribute {
	description := "Firewall rules for the service. **Required ports (must use exactly `[\"0.0.0.0/0\", \"::/0\"]`):** `22/tcp` (SSH), `4242/udp` (Nebula VPN)"

	if supportsCustomDomains {
		description += ", `80/tcp` (Let's Encrypt for custom domains)"
	}

	description += ". When `firewall_enabled` is `false`, set to `[]`."

	if len(templateRules) > 0 {
		description += " **Software default ports:** "
		var ruleParts []string
		for _, rule := range templateRules {
			ruleParts = append(ruleParts, fmt.Sprintf("`%s/%s`", rule.Port, rule.Protocol))
		}
		description += strings.Join(ruleParts, ", ") + "."
	}

	defaultRules := CreateUserRulesDefaultValue(templateRules, hasCustomFirewallPorts)
	defaultFirewallEnabled := hasCustomFirewallPorts

	return schema.SetNestedAttribute{
		MarkdownDescription: description,
		Optional:            true,
		Computed:            true,
		Default:             setdefault.StaticValue(defaultRules),
		Validators: []validator.Set{
			FirewallRulesConditional(defaultRules, defaultFirewallEnabled),
		},
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"type": schema.StringAttribute{
					MarkdownDescription: "The firewall rule type. Currently only `input` is supported.",
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf("input"),
					},
				},
				"port": schema.StringAttribute{
					MarkdownDescription: "Port number (`80`, `443`) or range (`8000-9000`).",
					Required:            true,
					Validators: []validator.String{
						validators.IsPortOrRange(),
					},
				},
				"protocol": schema.StringAttribute{
					MarkdownDescription: fmt.Sprintf("The protocol (`%s` or `%s`).", elestio.ServiceFirewallRuleProtocolTCP, elestio.ServiceFirewallRuleProtocolUDP),
					Required:            true,
					Validators: []validator.String{
						stringvalidator.OneOf(elestio.ServiceFirewallRuleProtocolTCP, elestio.ServiceFirewallRuleProtocolUDP),
					},
				},
				"targets": schema.SetAttribute{
					MarkdownDescription: "CIDR blocks allowed to access this port. Use `[\"0.0.0.0/0\", \"::/0\"]` to allow access from all IPv4/IPv6 addresses. **Note:** Required system ports (22/tcp/input, 4242/udp/input, and 80/tcp/input when using custom domains) must have exactly these two targets and cannot have additional target restrictions.",
					Required:            true,
					ElementType:         types.StringType,
				},
			},
		},
	}
}

// ToolRulesSchema returns the schema for firewall_tool_rules (read-only)
func ToolRulesSchema() schema.SetNestedAttribute {
	return schema.SetNestedAttribute{
		MarkdownDescription: "API-managed tool ports (VS Code, Terminal, File Explorer, etc.) that are not explicitly defined in `firewall_user_rules`." +
			" These ports are automatically managed by the API and will be preserved unless `firewall_remove_tool_ports` is set to `true`." +
			" If you want to manage a tool port yourself, include it explicitly in `firewall_user_rules`.",
		Computed: true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"name": schema.StringAttribute{
					MarkdownDescription: "The tool name.",
					Computed:            true,
				},
				"type": schema.StringAttribute{
					MarkdownDescription: "The firewall rule type.",
					Computed:            true,
				},
				"port": schema.StringAttribute{
					MarkdownDescription: "Port number.",
					Computed:            true,
				},
				"protocol": schema.StringAttribute{
					MarkdownDescription: "The protocol.",
					Computed:            true,
				},
				"targets": schema.SetAttribute{
					MarkdownDescription: "CIDR blocks allowed to access this port.",
					Computed:            true,
					ElementType:         types.StringType,
				},
			},
		},
	}
}
