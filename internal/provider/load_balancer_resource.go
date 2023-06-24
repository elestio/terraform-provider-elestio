package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/elestio/elestio-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &LoadBalancerResource{}
var _ resource.ResourceWithImportState = &LoadBalancerResource{}

func NewLoadBalancerResource() resource.Resource {
	return &LoadBalancerResource{}
}

type LoadBalancerResource struct {
	client *elestio.Client
}

type LoadBalancerResourceModel struct {
	Id        types.String `tfsdk:"id"`
	ProjectId types.String `tfsdk:"project_id"`
}

func (r *LoadBalancerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer"
}

func (r *LoadBalancerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Load balancer resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Load balancer identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Project identifier to which the load balancer will be attached",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			// TODO: add fields
		},
	}
}

func (r *LoadBalancerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *LoadBalancerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *LoadBalancerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Creating Load Balancer")

	loadBalancer, err := r.client.LoadBalancer.Create(elestio.CreateLoadBalancerRequest{
		ProjectID: plan.ProjectId.ValueString(),
		// TODO: add fields
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create load balancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created Load Balancer: "+loadBalancer.ID)

	loadBalancerModel := loadBalancerToResource(loadBalancer, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &loadBalancerModel)...)
}

func (r *LoadBalancerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *LoadBalancerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Getting Load Balancer: "+state.Id.ValueString())

	loadBalancer, err := r.client.LoadBalancer.Get(state.ProjectId.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read load balancer, got error: %s", err))
		return
	}
	loadBalancerModel := loadBalancerToResource(loadBalancer, state)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &loadBalancerModel)...)
}

func (r *LoadBalancerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state *LoadBalancerResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Updating Load Balancer: "+state.Id.ValueString())

	if plan.ProjectId.ValueString() != state.ProjectId.ValueString() {
		resp.Diagnostics.AddError("Error updating service", "Do not support project Id change")
		return
	}

	payload := elestio.UpdateLoadBalancerConfigRequest{
		// TODO: add fields
	}
	loadBalancer, err := r.client.LoadBalancer.UpdateConfig(state.ProjectId.ValueString(), state.Id.ValueString(), payload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update load balancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Updated Load Balancer: "+state.Id.ValueString())

	loadBalancerModel := loadBalancerToResource(loadBalancer, plan)
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &loadBalancerModel)...)
}

func (r *LoadBalancerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *LoadBalancerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Deleting Load Balancer: "+state.Id.ValueString())

	err := r.client.LoadBalancer.Delete(state.ProjectId.ValueString(), state.Id.ValueString(), true)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete load balancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Deleted Load Balancer: "+state.Id.ValueString())
}

func (r *LoadBalancerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: project_id,load_balancer_id. Got: %q", req.ID),
		)
		return
	}
	projectId := idParts[0]
	loadBalancerId := idParts[1]
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), projectId)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), loadBalancerId)...)
}

func loadBalancerToResource(l *elestio.LoadBalancer, state *LoadBalancerResourceModel) LoadBalancerResourceModel {
	return LoadBalancerResourceModel{
		Id:        types.StringValue(l.ID),
		ProjectId: types.StringValue(l.ProjectID),
		// TODO: add fields
	}
}
