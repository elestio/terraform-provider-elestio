package provider

import (
	"context"
	"os"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ provider.Provider = &ElestioProvider{}
)

type (
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
	resources := GetDeprecatedServiceResources()

	// Generate service resources from templates
	templateResources, err := GenerateServiceResourcesFromTemplates()
	if err != nil {
		panic(err)
	}

	resources = append(resources, templateResources...)
	return resources
}
