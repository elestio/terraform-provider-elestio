package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/elestio/elestio-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ProjectResource{}
	_ resource.ResourceWithImportState = &ProjectResource{}
)

type (
	ProjectResource struct {
		client *elestio.Client
	}

	ProjectResourceModel struct {
		Id              types.String `tfsdk:"id"`
		Name            types.String `tfsdk:"name"`
		Description     types.String `tfsdk:"description"`
		TechnicalEmails types.String `tfsdk:"technical_emails"`
		NetworkCIDR     types.String `tfsdk:"network_cidr"`
		CreationDate    types.String `tfsdk:"creation_date"`
		LastUpdated     types.String `tfsdk:"last_updated"`
	}
)

func NewProjectResource() resource.Resource {
	return &ProjectResource{}
}

func (r *ProjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *ProjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Project resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Project identifier.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Project name. Must be unique.",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Project description.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"technical_emails": schema.StringAttribute{
				MarkdownDescription: "Project technical emails.",
				Required:            true,
			},
			"network_cidr": schema.StringAttribute{
				MarkdownDescription: "Project network CIDR.",
				Computed:            true,
			},
			"creation_date": schema.StringAttribute{
				MarkdownDescription: "Project creation date.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ProjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*elestio.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *elestio.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ProjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.Project.Create(
		elestio.CreateProjectRequest{
			Name:            data.Name.ValueString(),
			Description:     data.Description.ValueString(),
			TechnicalEmails: data.TechnicalEmails.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Project",
			fmt.Sprintf("Unable to create project, got error: %s", err),
		)
		return
	}

	data.Id = types.StringValue(project.ID.String())
	UpdateTerraformDataWithElestioProject(data, project)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId := data.Id.ValueString()
	project, err := r.client.Project.Get(projectId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Project",
			fmt.Sprintf("Unable to read project, got error: %s", err),
		)
		return
	}

	UpdateTerraformDataWithElestioProject(data, project)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ProjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *ProjectResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	project, err := r.client.Project.Update(
		plan.Id.ValueString(),
		elestio.UpdateProjectRequest{
			Name:            plan.Name.ValueString(),
			Description:     plan.Description.ValueString(),
			TechnicalEmails: plan.TechnicalEmails.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Project",
			fmt.Sprintf("Unable to update project, got error: %s", err),
		)
		return
	}

	UpdateTerraformDataWithElestioProject(plan, project)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ProjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ProjectResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	serviceId := data.Id.ValueString()
	if err := r.client.Project.Delete(serviceId); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting  Project",
			fmt.Sprintf("Unable to delete project, got error: %s", err),
		)
		return
	}
}

func (r *ProjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func UpdateTerraformDataWithElestioProject(data *ProjectResourceModel, project *elestio.Project) {
	data.Name = types.StringValue(project.Name)
	data.Description = types.StringValue(project.Description)
	data.TechnicalEmails = types.StringValue(project.TechnicalEmails)
	data.NetworkCIDR = types.StringValue(project.NetworkCIDR)
	data.CreationDate = types.StringValue(project.CreationDate)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
}
