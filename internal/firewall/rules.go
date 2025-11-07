package firewall

import (
	"github.com/elestio/elestio-go-api-client/v2"
)

var toolPorts = map[string]string{
	"18345|tcp|input": "VS Code",
	"18374|tcp|input": "Open Terminal",
	"18346|tcp|input": "File Explorer",
	"18445|tcp|input": "Tail Logs",
	"18344|tcp|input": "Terminal",
}

// IsToolPort checks if a port is an Elestio management interface port
func IsToolPort(port, protocol, ruleType string) bool {
	_, exists := toolPorts[MakeFirewallRuleKey(port, protocol, ruleType)]
	return exists
}

// GetToolName returns the name of the tool for a given port
func GetToolName(port, protocol, ruleType string) string {
	return toolPorts[MakeFirewallRuleKey(port, protocol, ruleType)]
}

// GetRequiredSystemPorts returns required system ports (SSH and Nebula) as API rules
func GetRequiredSystemPorts() []elestio.ServiceFirewallRule {
	requiredPorts := GetSystemRequiredPorts()
	rules := make([]elestio.ServiceFirewallRule, len(requiredPorts))
	for i, port := range requiredPorts {
		rules[i] = elestio.ServiceFirewallRule{
			Port:     port.Port,
			Protocol: port.Protocol,
			Type:     port.Type,
			Targets:  port.Targets,
		}
	}
	return rules
}
