package provider

import (
	"context"
	_ "embed"
	"encoding/json"
	"os"
	"strings"

	"github.com/elestio/elestio-go-api-client"
	"github.com/elestio/terraform-provider-elestio/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &ElestioProvider{}

	//go:embed templates.json
	templatesListBytes []byte
)

type (
	TemplatesList struct {
		Templates []struct {
			ID                  int64  `json:"id"`
			Name                string `json:"title"`
			Category            string `json:"category"`
			Description         string `json:"description"`
			Logo                string `json:"mainImage"`
			DockerHubImage      string `json:"dockerhub_image"`
			DockerHubDefaultTag string `json:"dockerhub_default_tag"`
			FirewallPorts       string `json:"firewallPorts"`
		} `json:"templates"`
	}

	ElestioProvider struct {
		// version is set to the provider version on release, "dev" when the
		// provider is built and ran locally, and "test" when running acceptance
		// testing.
		version string
	}

	ElestioProviderModel struct {
		Email    types.String `tfsdk:"email"`
		APIToken types.String `tfsdk:"api_token"`
	}
)

func (p *ElestioProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "elestio"
	resp.Version = p.version
}

func (p *ElestioProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"email": schema.StringAttribute{
				MarkdownDescription: "Elestio email address." +
					" This is the email address with which you registered on the [Elestio website](https://dash.elest.io/).",
				Optional: true,
			},
			"api_token": schema.StringAttribute{
				MarkdownDescription: "Elestio API token." +
					" You can find this token in the [security settings](https://dash.elest.io/account/security) of your account.",
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *ElestioProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ElestioProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Email.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("email"),
			"Unknown Elestio API Email",
			"The provider cannot create the Elestio API client as there is an unknown configuration value for the Elestio API Email. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ELESTIO_EMAIL environment variable.",
		)
	}

	if data.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown Elestio API Token",
			"The provider cannot create the Elestio API client as there is an unknown configuration value for the Elestio API Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ELESTIO_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	email := os.Getenv("ELESTIO_EMAIL")
	apiToken := os.Getenv("ELESTIO_API_TOKEN")

	if !data.Email.IsNull() {
		email = data.Email.ValueString()
	}

	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	if email == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("email"),
			"Missing Elestio API Email",
			"The provider cannot create the Elestio API client as there is a missing or empty value for the Elestio API Email. "+
				"Set the host value in the configuration or use the ELESTIO_EMAIL environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing Elestio API Token",
			"The provider cannot create the Elestio API client as there is a missing or empty value for the Elestio API Token. "+
				"Set the host value in the configuration or use the ELESTIO_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client, err := elestio.NewClient(email, apiToken)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Elestio API Client",
			"An unexpected error occurred when creating the Elestio API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Elestio Client Error: "+err.Error(),
		)
		return
	}

	// Make the Elestio client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ElestioProvider) Resources(ctx context.Context) []func() resource.Resource {
	resources := []func() resource.Resource{
		NewProjectResource,
		NewLoadBalancerResource,
	}

	resources = append(resources, NewServiceResources()...)

	return resources
}

func (p *ElestioProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ElestioProvider{
			version: version,
		}
	}
}

func NewServiceResources() []func() resource.Resource {
	servicesResources := []func() resource.Resource{
		// Deprecated service resources
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:        0,
				ResourceName:      "service",
				DocumentationName: "Service",
				DeprecationMessage: "Use elestio_<SERVICE> resources instead. " +
					"This resource will be removed in the next major version of the provider.",
				Category: "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:        11,
				ResourceName:      "postgres",
				DocumentationName: "PostgreSQL",
				DeprecationMessage: "Use elestio_postgresql resource instead. " +
					"This resource will be removed in the next major version of the provider.",
				Category: "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         236,
				ResourceName:       "linux_desktop",
				DocumentationName:  "Linux-desktop",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         244,
				ResourceName:       "filerun",
				DocumentationName:  "FileRun",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         360,
				ResourceName:       "chaskiq",
				DocumentationName:  "Chaskiq",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         157,
				ResourceName:       "airbyte",
				DocumentationName:  "Airbyte",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         347,
				ResourceName:       "cal_com",
				DocumentationName:  "Cal.com",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         342,
				ResourceName:       "windmill",
				DocumentationName:  "Windmill",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         112,
				ResourceName:       "jupyter",
				DocumentationName:  "Jupyter",
				DeprecationMessage: "This resource was replaced by elestio_jupyter_notebook resource.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         185,
				ResourceName:       "opensourcetranslate",
				DocumentationName:  "OpenSourceTranslate",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
		func() resource.Resource {
			return NewServiceResource(&ServiceTemplate{
				TemplateId:         19,
				ResourceName:       "mongodb",
				DocumentationName:  "MongoDB",
				DeprecationMessage: "This resource is no more supported by Elestio and will be removed in the next major version of the provider.",
				Category:           "Deprecated",
			})
		},
	}

	// Unmarshal the bytes into the User struct
	var templatesList TemplatesList
	err := json.Unmarshal(templatesListBytes, &templatesList)
	if err != nil {
		panic(err)
	}

	for _, template := range templatesList.Templates {
		template := template // avoid iteration with same pointer

		// Skip full stack templates
		if template.Category == "Full Stack" {
			continue
		}

		var templateFirewallPorts []elestio.ServiceFirewallPort

		// Exemple of the template.FirewallPorts string: "80,443,22000,22000/udp"
		if len(template.FirewallPorts) > 0 {
			ports := strings.Split(template.FirewallPorts, ",")

			for _, p := range ports {
				var port string
				var protocol elestio.ServiceFirewallPortProtocol

				// Check if port is TCP or UDP
				if strings.Contains(p, "/udp") {
					port = strings.Split(p, "/")[0]
					protocol = elestio.ServiceFirewallPortProtocolUDP
				} else {
					port = p
					protocol = elestio.ServiceFirewallPortProtocolTCP
				}

				templateFirewallPorts = append(templateFirewallPorts, elestio.ServiceFirewallPort{
					// Remove whitespace if there is misstakes in the template
					Port:     strings.TrimSpace(port),
					Protocol: protocol,
				})
			}
		}

		servicesResources = append(
			servicesResources,
			func() resource.Resource {
				return NewServiceResource(&ServiceTemplate{
					TemplateId:        template.ID,
					ResourceName:      utils.CleanString(template.Name),
					DocumentationName: template.Name,
					Description:       template.Description,
					Logo:              template.Logo,
					DockerHubImage:    template.DockerHubImage,
					DefaultVersion:    func() string {
						if template.DockerHubDefaultTag == "" {
							return "latest"
						}
						return template.DockerHubDefaultTag
					}(),
					Category:          template.Category,
					FirewallPorts:     templateFirewallPorts,
				})
			},
		)
	}

	return servicesResources
}
