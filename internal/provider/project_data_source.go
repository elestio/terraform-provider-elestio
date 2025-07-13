package provider

import (
	"context"
	"fmt"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ProjectDataSource{}

type (
	ProjectDataSource struct {
		client *elestio.Client
	}

	ProjectDataSourceModel struct {
		Id              types.String `tfsdk:"id"`
		Name            types.String `tfsdk:"name"`
		Description     types.String `tfsdk:"description"`
		TechnicalEmails types.String `tfsdk:"technical_emails"` // deprecated
		TechnicalEmail  types.String `tfsdk:"technical_email"`
		NetworkCIDR     types.String `tfsdk:"network_cidr"`
		CreationDate    types.String `tfsdk:"creation_date"`
	}
)

func NewProjectDataSource() datasource.DataSource {
	return &ProjectDataSource{}
}

func (d *ProjectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *ProjectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Project data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Project description",
				Computed:            true,
			},
			"technical_emails": schema.StringAttribute{
				MarkdownDescription: "Project technical emails",
				DeprecationMessage:  "Use `technical_email` instead",
				Computed:            true,
			},
			"technical_email": schema.StringAttribute{
				MarkdownDescription: "Email address which will receive technical notifications",
				Computed:            true,
			},
			"network_cidr": schema.StringAttribute{
				MarkdownDescription: "Project network CIDR",
				Computed:            true,
			},
			"creation_date": schema.StringAttribute{
				MarkdownDescription: "Project creation date",
				Computed:            true,
			},
		},
	}
}

func (d *ProjectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*elestio.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *elestio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ProjectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ProjectDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := d.client.Project.Get(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Elestio Project",
			err.Error(),
		)
		return
	}

	state.Id = types.StringValue(project.ID.String())
	state.Name = types.StringValue(project.Name)
	state.Description = types.StringValue(project.Description)
	state.TechnicalEmails = types.StringValue(project.TechnicalEmail)
	state.TechnicalEmail = types.StringValue(project.TechnicalEmail)
	state.NetworkCIDR = types.StringValue(project.NetworkCIDR)
	state.CreationDate = types.StringValue(project.CreationDate)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
