package firewall

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	defaultTargetIPv4 = "0.0.0.0/0"
	defaultTargetIPv6 = "::/0"

	portProtocolSeparator = "|"
)

// RequiredPort defines a required firewall port with its specifications.
type RequiredPort struct {
	Port        string   // e.g., "22", "80", "4242"
	Protocol    string   // "tcp" or "udp"
	Type        string   // "input" or "output" (lowercase in Terraform, mapped to uppercase in API)
	Targets     []string // Required target CIDR blocks
	DisplayName string   // Human-readable description
}

// GetDefaultTargets returns default CIDR targets (all IPv4 and IPv6)
func GetDefaultTargets() []string {
	return []string{defaultTargetIPv4, defaultTargetIPv6}
}

// MakeFirewallRuleKey creates a unique identifier for a firewall rule (port|protocol|type)
func MakeFirewallRuleKey(port, protocol, ruleType string) string {
	return port + portProtocolSeparator + protocol + portProtocolSeparator + ruleType
}

// GetSystemRequiredPorts returns required system ports (SSH and Nebula)
func GetSystemRequiredPorts() []RequiredPort {
	return []RequiredPort{
		{
			Port:        "22",
			Protocol:    "tcp",
			Type:        "input",
			Targets:     GetDefaultTargets(),
			DisplayName: "22/tcp/input (SSH)",
		},
		{
			Port:        "4242",
			Protocol:    "udp",
			Type:        "input",
			Targets:     GetDefaultTargets(),
			DisplayName: "4242/udp/input (Nebula Global IP)",
		},
	}
}

// GetCustomDomainRequiredPorts returns required port 80 for Let's Encrypt SSL certificates
func GetCustomDomainRequiredPorts() []RequiredPort {
	return []RequiredPort{
		{
			Port:        "80",
			Protocol:    "tcp",
			Type:        "input",
			Targets:     GetDefaultTargets(),
			DisplayName: "80/tcp/input (HTTP for Let's Encrypt)",
		},
	}
}

// ValidateRequiredPorts validates that required ports are present with correct targets
func ValidateRequiredPorts(
	ctx context.Context,
	firewallUserRules types.Set,
	requiredPorts []RequiredPort,
	diags *diag.Diagnostics,
	attributePath path.Path,
	messageContext string,
) bool {
	if firewallUserRules.IsNull() || firewallUserRules.IsUnknown() {
		if len(requiredPorts) > 0 {
			var missingPorts []string
			for _, reqPort := range requiredPorts {
				portMsg := fmt.Sprintf("%s must be added to firewall_user_rules with targets: %v", reqPort.DisplayName, reqPort.Targets)
				if messageContext != "" {
					portMsg = fmt.Sprintf("%s %s", messageContext, portMsg)
				}
				missingPorts = append(missingPorts, portMsg)
			}
			message := buildErrorMessage(missingPorts)
			diags.AddAttributeError(attributePath, "Invalid Firewall Configuration", message)
		}
		return false
	}

	var userRules []FirewallRuleModel
	d := firewallUserRules.ElementsAs(ctx, &userRules, true)
	diags.Append(d...)
	if diags.HasError() {
		return false
	}

	// Track which required ports are present and their targets
	foundPorts := make(map[string][]string)
	requiredPortsMap := make(map[string]RequiredPort)

	for _, reqPort := range requiredPorts {
		key := MakeFirewallRuleKey(reqPort.Port, reqPort.Protocol, reqPort.Type)
		requiredPortsMap[key] = reqPort
	}

	// Scan user rules to find required ports
	for _, rule := range userRules {
		port := rule.Port.ValueString()
		protocol := rule.Protocol.ValueString()
		ruleType := rule.Type.ValueString()

		key := MakeFirewallRuleKey(port, protocol, ruleType)

		if _, isRequired := requiredPortsMap[key]; isRequired {
			var targets []string
			d := rule.Targets.ElementsAs(ctx, &targets, true)
			diags.Append(d...)
			if diags.HasError() {
				return false
			}
			foundPorts[key] = targets
		}
	}

	// Check for missing required ports and validate targets
	var errors []string
	for key, reqPort := range requiredPortsMap {
		targets, found := foundPorts[key]
		if !found {
			portMsg := fmt.Sprintf("%s must be added to firewall_user_rules with targets: %v", reqPort.DisplayName, reqPort.Targets)
			if messageContext != "" {
				portMsg = fmt.Sprintf("%s %s", messageContext, portMsg)
			}
			errors = append(errors, portMsg)
			continue
		}

		// Check if all required targets are present
		targetMap := make(map[string]bool)
		for _, target := range targets {
			targetMap[target] = true
		}

		var missingTargets []string
		for _, requiredTarget := range reqPort.Targets {
			if !targetMap[requiredTarget] {
				missingTargets = append(missingTargets, requiredTarget)
			}
		}

		if len(missingTargets) > 0 {
			errors = append(errors, fmt.Sprintf("%s is present but missing required targets: %v", reqPort.DisplayName, missingTargets))
		}

		// Check if there are extra targets (strict validation for system ports)
		if len(targets) != len(reqPort.Targets) {
			var extraTargets []string
			requiredTargetMap := make(map[string]bool)
			for _, requiredTarget := range reqPort.Targets {
				requiredTargetMap[requiredTarget] = true
			}
			for _, target := range targets {
				if !requiredTargetMap[target] {
					extraTargets = append(extraTargets, target)
				}
			}
			if len(extraTargets) > 0 {
				errors = append(errors, fmt.Sprintf("%s must have exactly targets %v, but found additional targets: %v", reqPort.DisplayName, reqPort.Targets, extraTargets))
			}
		}
	}

	if len(errors) > 0 {
		message := buildErrorMessage(errors)
		diags.AddAttributeError(attributePath, "Invalid Firewall Configuration", message)
		return false
	}

	return true
}

func buildErrorMessage(errors []string) string {
	if len(errors) == 1 {
		return errors[0]
	}
	message := "Missing required ports:\n"
	for i, err := range errors {
		message += fmt.Sprintf("  - %s", err)
		if i < len(errors)-1 {
			message += "\n"
		}
	}
	return message
}

// FirewallRuleModel represents a firewall rule in Terraform
type FirewallRuleModel struct {
	Port     types.String `tfsdk:"port"`
	Protocol types.String `tfsdk:"protocol"`
	Type     types.String `tfsdk:"type"`
	Targets  types.Set    `tfsdk:"targets"`
}
