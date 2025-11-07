package provider

import (
	_ "embed"
	"encoding/json"
	"strings"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/firewall"
	"github.com/elestio/terraform-provider-elestio/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	//go:embed templates.json
	templatesListBytes []byte
)

type (
	// TemplatesList represents the structure of templates.json
	TemplatesList struct {
		Templates []RawTemplate `json:"templates"`
	}

	// RawTemplate represents a template as it appears in templates.json
	RawTemplate struct {
		ID                  int64  `json:"id"`
		Name                string `json:"title"`
		Category            string `json:"category"`
		Description         string `json:"description"`
		Logo                string `json:"mainImage"`
		DockerHubImage      string `json:"dockerhub_image"`
		DockerHubDefaultTag string `json:"dockerhub_default_tag"`
		FirewallPorts       string `json:"firewallPorts"`
	}

	// DeprecatedResource represents a deprecated service resource configuration
	DeprecatedResource struct {
		TemplateId         int64
		ResourceName       string
		DocumentationName  string
		DeprecationMessage string
	}
)

// getDeprecatedResourcesConfig returns the configuration for all deprecated resources
func getDeprecatedResourcesConfig() []DeprecatedResource {
	return []DeprecatedResource{
		{
			TemplateId:        0,
			ResourceName:      "service",
			DocumentationName: "Service",
			DeprecationMessage: "Use elestio_<SERVICE> resources instead. " +
				"This resource will be removed in the next major version of the provider.",
		},
		{
			TemplateId:        11,
			ResourceName:      "postgres",
			DocumentationName: "PostgreSQL",
			DeprecationMessage: "Use elestio_postgresql resource instead. " +
				"This resource will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         236,
			ResourceName:       "linux_desktop",
			DocumentationName:  "Linux-desktop",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         244,
			ResourceName:       "filerun",
			DocumentationName:  "FileRun",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         360,
			ResourceName:       "chaskiq",
			DocumentationName:  "Chaskiq",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         157,
			ResourceName:       "airbyte",
			DocumentationName:  "Airbyte",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         347,
			ResourceName:       "cal_com",
			DocumentationName:  "Cal.com",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         342,
			ResourceName:       "windmill",
			DocumentationName:  "Windmill",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         112,
			ResourceName:       "jupyter",
			DocumentationName:  "Jupyter",
			DeprecationMessage: "This resource was replaced by elestio_jupyter_notebook resource.",
		},
		{
			TemplateId:         185,
			ResourceName:       "opensourcetranslate",
			DocumentationName:  "OpenSourceTranslate",
			DeprecationMessage: "This resource was replaced by elestio_osstranslate resource.",
		},
		{
			TemplateId:         19,
			ResourceName:       "mongodb",
			DocumentationName:  "MongoDB",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         185, // Note: Duplicate ID with opensourcetranslate
			ResourceName:       "libretranslate",
			DocumentationName:  "LibreTranslate",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         401,
			ResourceName:       "gophish",
			DocumentationName:  "Gophish",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
		{
			TemplateId:         375,
			ResourceName:       "surrealdb",
			DocumentationName:  "SurrealDB",
			DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
		},
	}
}

// LoadTemplatesList unmarshals the embedded templates.json file
func LoadTemplatesList() (*TemplatesList, error) {
	var templatesList TemplatesList
	err := json.Unmarshal(templatesListBytes, &templatesList)
	if err != nil {
		return nil, err
	}
	return &templatesList, nil
}

// ParseFirewallRules parses firewall ports string and returns firewall rules (template-specific only, without required system ports)
func ParseFirewallRules(firewallPorts string) ([]elestio.ServiceFirewallRule, bool) {
	var templateFirewallRules []elestio.ServiceFirewallRule
	hasCustomFirewallPorts := len(firewallPorts) > 0

	// Example of the firewallPorts string: "80,443,22000,22000/udp"
	if hasCustomFirewallPorts {
		ports := strings.Split(firewallPorts, ",")

		for _, p := range ports {
			rule := parseFirewallPort(strings.TrimSpace(p))
			templateFirewallRules = append(templateFirewallRules, rule)
		}
	}

	// Required system ports (22/TCP for SSH, 4242/UDP for Nebula) are added automatically by the API
	return templateFirewallRules, hasCustomFirewallPorts
}

// parseFirewallPort parses a single port string and returns a firewall rule
func parseFirewallPort(port string) elestio.ServiceFirewallRule {
	var portNum, protocol string

	// Check if port is TCP or UDP
	if strings.Contains(port, "/udp") {
		portNum = strings.Split(port, "/")[0]
		protocol = elestio.ServiceFirewallRuleProtocolUDP
	} else {
		portNum = port
		protocol = elestio.ServiceFirewallRuleProtocolTCP
	}

	return createFirewallRule(portNum, protocol)
}

func createFirewallRule(port, protocol string) elestio.ServiceFirewallRule {
	return elestio.ServiceFirewallRule{
		Port:     port,
		Protocol: protocol,
		Type:     elestio.ServiceFirewallRuleTypeInput,
		Targets:  firewall.GetDefaultTargets(),
	}
}

// CreateServiceTemplateFromRaw creates a ServiceTemplate from raw template data
func CreateServiceTemplateFromRaw(template RawTemplate) *ServiceTemplate {
	firewallRules, hasCustomFirewallPorts := ParseFirewallRules(template.FirewallPorts)

	defaultVersion := template.DockerHubDefaultTag
	if defaultVersion == "" {
		defaultVersion = "latest"
	}

	return &ServiceTemplate{
		TemplateId:             template.ID,
		ResourceName:           utils.CleanString(template.Name),
		DocumentationName:      template.Name,
		Description:            template.Description,
		Logo:                   template.Logo,
		DockerHubImage:         template.DockerHubImage,
		DefaultVersion:         defaultVersion,
		Category:               template.Category,
		FirewallRules:          firewallRules,
		HasCustomFirewallPorts: hasCustomFirewallPorts,
	}
}

// GetDeprecatedServiceResources returns the list of deprecated service resources
func GetDeprecatedServiceResources() []func() resource.Resource {
	deprecatedConfigs := getDeprecatedResourcesConfig()
	resources := make([]func() resource.Resource, 0, len(deprecatedConfigs))

	for _, config := range deprecatedConfigs {
		// Capture the config in the closure to avoid variable capture issues
		config := config
		resources = append(resources, func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:             config.TemplateId,
				ResourceName:           config.ResourceName,
				DocumentationName:      config.DocumentationName,
				DeprecationMessage:     config.DeprecationMessage,
				Category:               "Deprecated",
				HasCustomFirewallPorts: false,
			})
		})
	}

	return resources
}

// GenerateServiceResourcesFromTemplates creates service resources from the templates.json file
func GenerateServiceResourcesFromTemplates() ([]func() resource.Resource, error) {
	templatesList, err := LoadTemplatesList()
	if err != nil {
		return nil, err
	}

	var servicesResources []func() resource.Resource

	for _, template := range templatesList.Templates {
		// Skip full stack templates
		if template.Category == "Full Stack" {
			continue
		}

		// Create the service template outside the closure to avoid capture issues
		serviceTemplate := CreateServiceTemplateFromRaw(template)

		servicesResources = append(servicesResources, func() resource.Resource {
			return NewServiceResource(serviceTemplate)
		})
	}

	return servicesResources, nil
}
