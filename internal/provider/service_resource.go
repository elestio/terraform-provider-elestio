package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elestio/elestio-go-api-client"
	"github.com/elestio/terraform-provider-elestio/internal/modifiers"
	"github.com/elestio/terraform-provider-elestio/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sdk_resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var (
	_ resource.Resource                   = &ServiceResource{}
	_ resource.ResourceWithSchema         = &ServiceResource{}
	_ resource.ResourceWithValidateConfig = &ServiceResource{}
	_ resource.ResourceWithConfigure      = &ServiceResource{}
	_ resource.ResourceWithImportState    = &ServiceResource{}
)

type (
	ServiceTemplate struct {
		TemplateId         int64
		ResourceName       string
		DocumentationName  string
		DeprecationMessage string
	}

	ServiceResource struct {
		client *elestio.Client
		*ServiceTemplate
	}

	ServiceResourceAdminModel struct {
		URL      types.String `tfsdk:"url"`
		User     types.String `tfsdk:"user"`
		Password types.String `tfsdk:"password"`
	}

	ServiceResourceDatabaseAdminModel struct {
		Host     types.String `tfsdk:"host"`
		Port     types.String `tfsdk:"port"`
		User     types.String `tfsdk:"user"`
		Password types.String `tfsdk:"password"`
		Command  types.String `tfsdk:"command"`
	}

	ServiceResourceModel struct {
		Id                                          types.String `tfsdk:"id"`
		ProjectID                                   types.String `tfsdk:"project_id"`
		ServerName                                  types.String `tfsdk:"server_name"`
		ServerType                                  types.String `tfsdk:"server_type"`
		TemplateId                                  types.Int64  `tfsdk:"template_id"`
		Version                                     types.String `tfsdk:"version"`
		ProviderName                                types.String `tfsdk:"provider_name"`
		Datacenter                                  types.String `tfsdk:"datacenter"`
		SupportLevel                                types.String `tfsdk:"support_level"`
		AdminEmail                                  types.String `tfsdk:"admin_email"`
		Category                                    types.String `tfsdk:"category"`
		Status                                      types.String `tfsdk:"status"`
		DeploymentStatus                            types.String `tfsdk:"deployment_status"`
		DeploymentStartedAt                         types.String `tfsdk:"deployment_started_at"`
		DeploymentEndedAt                           types.String `tfsdk:"deployment_ended_at"`
		CreatorName                                 types.String `tfsdk:"creator_name"`
		CreatedAt                                   types.String `tfsdk:"created_at"`
		IPV4                                        types.String `tfsdk:"ipv4"`
		IPV6                                        types.String `tfsdk:"ipv6"`
		CNAME                                       types.String `tfsdk:"cname"`
		Country                                     types.String `tfsdk:"country"`
		City                                        types.String `tfsdk:"city"`
		AdminUser                                   types.String `tfsdk:"admin_user"`
		RootAppPath                                 types.String `tfsdk:"root_app_path"`
		Env                                         types.Map    `tfsdk:"env"`
		Admin                                       types.Object `tfsdk:"admin"`
		DatabaseAdmin                               types.Object `tfsdk:"database_admin"`
		GlobalIP                                    types.String `tfsdk:"global_ip"`
		TrafficOutgoing                             types.Int64  `tfsdk:"traffic_outgoing"`
		TrafficIncoming                             types.Int64  `tfsdk:"traffic_incoming"`
		TrafficIncluded                             types.Int64  `tfsdk:"traffic_included"`
		Cores                                       types.Int64  `tfsdk:"cores"`
		RAMSizeGB                                   types.String `tfsdk:"ram_size_gb"`
		StorageSizeGB                               types.Int64  `tfsdk:"storage_size_gb"`
		PricePerHour                                types.String `tfsdk:"price_per_hour"`
		AppAutoUpdatesEnabled                       types.Bool   `tfsdk:"app_auto_updates_enabled"`
		AppAutoUpdatesDayOfWeek                     types.Int64  `tfsdk:"app_auto_updates_day_of_week"`
		AppAutoUpdatesHour                          types.Int64  `tfsdk:"app_auto_updates_hour"`
		AppAutoUpdatesMinute                        types.Int64  `tfsdk:"app_auto_updates_minute"`
		SystemAutoUpdatesEnabled                    types.Bool   `tfsdk:"system_auto_updates_enabled"`
		SystemAutoUpdatesSecurityPatchesOnlyEnabled types.Bool   `tfsdk:"system_auto_updates_security_patches_only_enabled"`
		SystemAutoUpdatesRebootDayOfWeek            types.Int64  `tfsdk:"system_auto_updates_reboot_day_of_week"`
		SystemAutoUpdatesRebootHour                 types.Int64  `tfsdk:"system_auto_updates_reboot_hour"`
		SystemAutoUpdatesRebootMinute               types.Int64  `tfsdk:"system_auto_updates_reboot_minute"`
		BackupsEnabled                              types.Bool   `tfsdk:"backups_enabled"`
		RemoteBackupsEnabled                        types.Bool   `tfsdk:"remote_backups_enabled"`
		ExternalBackupsEnabled                      types.Bool   `tfsdk:"external_backups_enabled"`
		ExternalBackupsUpdateDayOfWeek              types.Int64  `tfsdk:"external_backups_update_day_of_week"`
		ExternalBackupsUpdateHour                   types.Int64  `tfsdk:"external_backups_update_hour"`
		ExternalBackupsUpdateMinute                 types.Int64  `tfsdk:"external_backups_update_minute"`
		ExternalBackupsUpdateType                   types.String `tfsdk:"external_backups_update_type"`
		ExternalBackupsRetainDayOfWeek              types.Int64  `tfsdk:"external_backups_retain_day_of_week"`
		FirewallEnabled                             types.Bool   `tfsdk:"firewall_enabled"`
		FirewallId                                  types.String `tfsdk:"firewall_id"`
		FirewallPorts                               types.String `tfsdk:"firewall_ports"`
		AlertsEnabled                               types.Bool   `tfsdk:"alerts_enabled"`
		LastUpdated                                 types.String `tfsdk:"last_updated"`
	}
)

func NewServiceResource(template *ServiceTemplate) resource.Resource {
	return &ServiceResource{
		ServiceTemplate: template,
	}
}

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.ResourceName
}

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: r.DocumentationName + " resource",
		DeprecationMessage:  r.DeprecationMessage,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Service identifier.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.StringAttribute{
				MarkdownDescription: "Identifier of the project the service will be created.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_name": schema.StringAttribute{
				MarkdownDescription: "Service server name. Must be unique within the project.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Service server type. You can only upgrade it, not downgrade.",
				Required:            true,
			},
			"template_id": schema.Int64Attribute{
				MarkdownDescription: "Service template identifier.",
				Required:            r.TemplateId == 0,
				Computed:            r.TemplateId != 0,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Service software version.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						func(ctx context.Context, modifier planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
							// PostgreSQL = 11
							if r.TemplateId == 11 {
								resp.RequiresReplace = true
								return
							}

							resp.RequiresReplace = false
						},
						"Requires replace if you want to upgrade version.",
						"Requires replace if you want to upgrade version.",
					),
				},
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Service provider name.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"datacenter": schema.StringAttribute{
				MarkdownDescription: "Service datacenter.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"support_level": schema.StringAttribute{
				MarkdownDescription: "Service support level.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("level1", "level2", "level3"),
				},
			},
			"admin_email": schema.StringAttribute{
				MarkdownDescription: "Service admin email.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Service category.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Service status.",
				Computed:            true,
			},
			"deployment_status": schema.StringAttribute{
				MarkdownDescription: "Service deployement status.",
				Computed:            true,
			},
			"deployment_started_at": schema.StringAttribute{
				MarkdownDescription: "Service deployment startedAt date.",
				Computed:            true,
			},
			"deployment_ended_at": schema.StringAttribute{
				MarkdownDescription: "Service deployment endedAt date.",
				Computed:            true,
			},
			"creator_name": schema.StringAttribute{
				MarkdownDescription: "Service creator name.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Service creation date.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "Service IPv4.",
				Computed:            true,
			},
			"ipv6": schema.StringAttribute{
				MarkdownDescription: "Service IPv6.",
				Computed:            true,
			},
			"cname": schema.StringAttribute{
				MarkdownDescription: "Service CNAME.",
				Computed:            true,
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "Service country.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "Service city.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_user": schema.StringAttribute{
				MarkdownDescription: "Service admin user.",
				Computed:            true,
			},
			"root_app_path": schema.StringAttribute{
				MarkdownDescription: "Service root app path.",
				Computed:            true,
			},
			"env": schema.MapAttribute{
				MarkdownDescription: "Service environment variables.",
				ElementType:         types.StringType,
				Computed:            true,
				Sensitive:           true,
			},
			"admin": schema.SingleNestedAttribute{
				MarkdownDescription: "Service admin.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						MarkdownDescription: "Service admin URL.",
						Computed:            true,
					},
					"user": schema.StringAttribute{
						MarkdownDescription: "Service admin user.",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Service admin password.",
						Computed:            true,
						Sensitive:           true,
					},
				},
			},
			"database_admin": schema.SingleNestedAttribute{
				MarkdownDescription: "Service database admin.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						MarkdownDescription: "Service database admin host.",
						Computed:            true,
					},
					"port": schema.StringAttribute{
						MarkdownDescription: "Service database admin port.",
						Computed:            true,
					},
					"user": schema.StringAttribute{
						MarkdownDescription: "Service database admin user.",
						Computed:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Service database admin password.",
						Computed:            true,
						Sensitive:           true,
					},
					"command": schema.StringAttribute{
						MarkdownDescription: "Service database admin command.",
						Computed:            true,
						Sensitive:           true,
					},
				},
			},
			"global_ip": schema.StringAttribute{
				MarkdownDescription: "Service global IP.",
				Computed:            true,
			},
			"traffic_outgoing": schema.Int64Attribute{
				MarkdownDescription: "Service traffic outgoing.",
				Computed:            true,
			},
			"traffic_incoming": schema.Int64Attribute{
				MarkdownDescription: "Service traffic incoming.",
				Computed:            true,
			},
			"traffic_included": schema.Int64Attribute{
				MarkdownDescription: "Service traffic included.",
				Computed:            true,
			},
			"cores": schema.Int64Attribute{
				MarkdownDescription: "Service cores.",
				Computed:            true,
			},
			"ram_size_gb": schema.StringAttribute{
				MarkdownDescription: "Service ram size.",
				Computed:            true,
			},
			"storage_size_gb": schema.Int64Attribute{
				MarkdownDescription: "Service storage size.",
				Computed:            true,
			},
			"price_per_hour": schema.StringAttribute{
				MarkdownDescription: "Service price per hour.",
				Computed:            true,
			},
			"app_auto_updates_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service app auto update state. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(true),
				},
			},
			"app_auto_updates_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service app auto update day of week.",
				Computed:            true,
			},
			"app_auto_updates_hour": schema.Int64Attribute{
				MarkdownDescription: "Service app auto update hour.",
				Computed:            true,
			},
			"app_auto_updates_minute": schema.Int64Attribute{
				MarkdownDescription: "Service app auto update minute.",
				Computed:            true,
			},
			"system_auto_updates_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service system auto update state. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(true),
				},
			},
			"system_auto_updates_security_patches_only_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service system auto update security patches only state. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(false),
				},
			},
			"system_auto_updates_reboot_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service system auto update reboot day of week.",
				Computed:            true,
			},
			"system_auto_updates_reboot_hour": schema.Int64Attribute{
				MarkdownDescription: "Service system auto update reboot hour.",
				Computed:            true,
			},
			"system_auto_updates_reboot_minute": schema.Int64Attribute{
				MarkdownDescription: "Service system auto update reboot minute.",
				Computed:            true,
			},
			"backups_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service backups state. Requires a support_level higher than `level1`. **Default** `false`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(false),
				},
			},
			"remote_backups_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service remote backups state. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(true),
				},
			},
			"external_backups_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service external backups state. **Default** `false`.",
				Computed:            true,
				// TODO: Handle external backups with s3 config
			},
			"external_backups_update_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service external backups update day.",
				Computed:            true,
			},
			"external_backups_update_hour": schema.Int64Attribute{
				MarkdownDescription: "Service external backups update hour.",
				Computed:            true,
			},
			"external_backups_update_minute": schema.Int64Attribute{
				MarkdownDescription: "Service external backups update minute.",
				Computed:            true,
			},
			"external_backups_update_type": schema.StringAttribute{
				MarkdownDescription: "Service external backups update type.",
				Computed:            true,
			},
			"external_backups_retain_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service external backups retain day.",
				Computed:            true,
			},
			"firewall_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service firewall state. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(true),
				},
			},
			"firewall_id": schema.StringAttribute{
				MarkdownDescription: "Service firewall id.",
				Computed:            true,
			},
			"firewall_ports": schema.StringAttribute{
				MarkdownDescription: "Service firewall ports.",
				Computed:            true,
			},
			"alerts_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service alerts state. **Default** `true`.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					modifiers.BoolDefault(true),
				},
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *ServiceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data ServiceResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.BackupsEnabled.ValueBool() && data.SupportLevel.ValueString() == "level1" {
		resp.Diagnostics.AddAttributeError(
			path.Root("backups_enabled"),
			"Invalid Attribute Configuration",
			"The backups are available only for a support level higher than level1. "+
				"You must upgrade support_level to enable backups_enabled.",
		)
		return
	}

	if data.SystemAutoUpdatesSecurityPatchesOnlyEnabled.ValueBool() && !data.SystemAutoUpdatesEnabled.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("system_auto_updates_security_patches_only_enabled"),
			"Invalid Attribute Configuration",
			"The system_auto_updates_security_patches_only_enabled can be enabled only if system_auto_updates_enabled is enabled.",
		)
		return
	}
}

func (r *ServiceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If no template is provided in the provider configuration
	// use the one provided by the user.
	var templateId int64
	if r.TemplateId != 0 {
		templateId = r.TemplateId
	} else {
		templateId = data.TemplateId.ValueInt64()
	}

	// The service will be created but we should wait
	// for it to be fully deployed.
	serviceCreating, err := r.client.Service.Create(
		elestio.CreateServiceRequest{
			ProjectID:    data.ProjectID.ValueString(),
			ServerName:   data.ServerName.ValueString(),
			ServerType:   data.ServerType.ValueString(),
			TemplateID:   templateId,
			Version:      data.Version.ValueString(),
			ProviderName: data.ProviderName.ValueString(),
			Datacenter:   data.Datacenter.ValueString(),
			SupportLevel: data.SupportLevel.ValueString(),
			AdminEmail:   data.AdminEmail.ValueString(),
		},
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Service",
			fmt.Sprintf("Unable to start service creation, got error: %s", err),
		)
		return
	}

	serviceCreated, err := r.waitServiceCreate(ctx, serviceCreating)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Service",
			fmt.Sprintf("Unable to wait service creation, got error: %s", err),
		)
		return
	}

	data.Id = types.StringValue(serviceCreated.ID)

	// Update some fields that are not available in the create request.
	serviceUpdated, err := r.updateElestioService(ctx, serviceCreated, data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Service",
			fmt.Sprintf("Unable after create to update fields that are not available in the create request, got error: %s", err),
		)
		return
	}

	convertElestioToTerraformFormat(data, serviceUpdated, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId, serviceId := data.ProjectID.ValueString(), data.Id.ValueString()
	service, err := r.client.Service.Get(projectId, serviceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Service",
			fmt.Sprintf("Unable to read service, got error: %s", err),
		)
		return
	}

	convertElestioToTerraformFormat(data, service, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan *ServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId, serviceId := state.ProjectID.ValueString(), state.Id.ValueString()
	service, err := r.client.Service.Get(projectId, serviceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Service",
			fmt.Sprintf("Unable to get service, got error: %s", err),
		)
		return
	}

	updatedService, err := r.updateElestioService(ctx, service, plan, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Service",
			fmt.Sprintf("Unable to update service, got error: %s", err),
		)
		return
	}

	convertElestioToTerraformFormat(plan, updatedService, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId, serviceId := state.ProjectID.ValueString(), state.Id.ValueString()
	service, err := r.client.Service.Get(projectId, serviceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Service",
			fmt.Sprintf("Unable to get service, got error: %s", err),
		)
		return
	}

	if err := r.client.Service.Delete(service.ProjectID, service.ID); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Service",
			fmt.Sprintf("Unable to start service deletion, got error: %s", err),
		)
		return
	}

	if err := r.waitServiceDelete(ctx, service); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Service",
			fmt.Sprintf("Unable to wait service deletion, got error: %s", err),
		)
		return
	}
}

func (r *ServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: project_id,service_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *ServiceResource) updateElestioService(ctx context.Context, service *elestio.Service, plan *ServiceResourceModel, diags *diag.Diagnostics) (*elestio.Service, error) {
	state := &ServiceResourceModel{}
	state.Id = types.StringValue(service.ID)
	convertElestioToTerraformFormat(state, service, diags)

	// Server type update should be done first, because it requires to stop the service
	if !state.ServerType.Equal(plan.ServerType) {
		if err := r.client.Service.UpdateServerType(service.ID, plan.ServerType.ValueString(), service.ProviderName, service.Datacenter); err != nil {
			return nil, fmt.Errorf("failed to update serverType: %s", err)
		}
		r.waitServerTypeUpdate(ctx, service, plan.ServerType.ValueString())
	}

	if !state.Version.Equal(plan.Version) {
		if err := r.client.Service.UpdateVersion(service.ID, plan.Version.ValueString()); err != nil {
			return nil, fmt.Errorf("failed to update version: %s", err)
		}
	}

	if !state.AppAutoUpdatesEnabled.Equal(plan.AppAutoUpdatesEnabled) {
		if plan.AppAutoUpdatesEnabled.ValueBool() {
			if err := r.client.Service.EnableAppAutoUpdates(service.ID); err != nil {
				return nil, fmt.Errorf("failed to enable appAutoUpdates: %s", err)
			}
		} else {
			if err := r.client.Service.DisableAppAutoUpdates(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable appAutoUpdates: %s", err)
			}
		}
	}

	if !state.SystemAutoUpdatesEnabled.Equal(plan.SystemAutoUpdatesEnabled) {
		if plan.SystemAutoUpdatesEnabled.ValueBool() {
			if err := r.client.Service.EnableSystemAutoUpdates(service.ID, plan.SystemAutoUpdatesSecurityPatchesOnlyEnabled.ValueBool()); err != nil {
				return nil, fmt.Errorf("failed to enable systemAutoUpdates: %s", err)
			}
		} else {
			if err := r.client.Service.DisableSystemAutoUpdates(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable systemAutoUpdates: %s", err)
			}
		}
	}

	if !state.SystemAutoUpdatesSecurityPatchesOnlyEnabled.Equal(plan.SystemAutoUpdatesSecurityPatchesOnlyEnabled) {
		if err := r.client.Service.EnableSystemAutoUpdates(service.ID, plan.SystemAutoUpdatesSecurityPatchesOnlyEnabled.ValueBool()); err != nil {
			return nil, fmt.Errorf("failed to enable systemAutoUpdates: %s", err)
		}
	}

	if !state.AppAutoUpdatesEnabled.Equal(plan.AppAutoUpdatesEnabled) {
		if plan.AppAutoUpdatesEnabled.ValueBool() {
			if err := r.client.Service.EnableAppAutoUpdates(service.ID); err != nil {
				return nil, fmt.Errorf("failed to enable appAutoUpdates: %s", err)
			}
		} else {
			if err := r.client.Service.DisableAppAutoUpdates(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable appAutoUpdates: %s", err)
			}
		}
	}

	if !state.BackupsEnabled.Equal(plan.BackupsEnabled) {
		if plan.BackupsEnabled.ValueBool() {
			if err := r.client.Service.EnableBackups(service.ID); err != nil {
				return nil, fmt.Errorf("failed to enable backups: %s", err)
			}
		} else {
			if err := r.client.Service.DisableBackups(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable backups: %s", err)
			}
		}
	}

	if !state.RemoteBackupsEnabled.Equal(plan.RemoteBackupsEnabled) {
		if plan.RemoteBackupsEnabled.ValueBool() {
			if err := r.client.Service.EnableRemoteBackups(service.ID); err != nil {
				return nil, fmt.Errorf("failed to enable remoteBackups: %s", err)
			}
		} else {
			if err := r.client.Service.DisableRemoteBackups(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable remoteBackups: %s", err)
			}
		}
	}

	if !state.FirewallEnabled.Equal(plan.FirewallEnabled) {
		if plan.FirewallEnabled.ValueBool() {
			if err := r.client.Service.EnableFirewall(service.ID); err != nil {
				return nil, fmt.Errorf("failed to enable firewall: %s", err)
			}
		} else {
			if err := r.client.Service.DisableFirewall(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable firewall: %s", err)
			}
		}
	}

	if !state.AlertsEnabled.Equal(plan.AlertsEnabled) {
		if plan.AlertsEnabled.ValueBool() {
			if err := r.client.Service.EnableAlerts(service.ID); err != nil {
				return nil, fmt.Errorf("failed to enable alerts: %s", err)
			}
		} else {
			if err := r.client.Service.DisableAlerts(service.ID); err != nil {
				return nil, fmt.Errorf("failed to disable alerts: %s", err)
			}
		}
	}

	service, err := r.client.Service.Get(service.ProjectID, service.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %s", err)
	}

	return service, nil
}

func convertElestioToTerraformFormat(data *ServiceResourceModel, service *elestio.Service, diags *diag.Diagnostics) {
	data.ProjectID = types.StringValue(service.ProjectID)
	data.ServerName = types.StringValue(service.ServerName)
	data.ServerType = types.StringValue(service.ServerType)
	data.TemplateId = types.Int64Value(service.TemplateID)
	data.Version = types.StringValue(service.Version)
	data.ProviderName = types.StringValue(service.ProviderName)
	data.Datacenter = types.StringValue(service.Datacenter)
	data.SupportLevel = types.StringValue(service.SupportLevel)
	data.AdminEmail = types.StringValue(service.AdminEmail)
	data.Category = types.StringValue(service.Category)
	data.Status = types.StringValue(service.Status)
	data.DeploymentStatus = types.StringValue(service.DeploymentStatus)
	data.DeploymentStartedAt = types.StringValue(service.DeploymentStartedAt)
	data.DeploymentEndedAt = types.StringValue(service.DeploymentEndedAt)
	data.CreatorName = types.StringValue(service.CreatorName)
	data.CreatedAt = types.StringValue(service.CreatedAt)
	data.IPV4 = types.StringValue(service.IPV4)
	data.IPV6 = types.StringValue(service.IPV6)
	data.CNAME = types.StringValue(service.CNAME)
	data.Country = types.StringValue(service.Country)
	data.City = types.StringValue(service.City)
	data.AdminUser = types.StringValue(service.AdminUser)
	data.RootAppPath = types.StringValue(service.RootAppPath)
	data.Env = utils.MapStringToMapType(service.Env, diags)
	data.Admin = utils.ObjectValue(
		map[string]attr.Type{
			"url":      types.StringType,
			"user":     types.StringType,
			"password": types.StringType,
		},
		map[string]attr.Value{
			"url":      types.StringValue(service.Admin.URL),
			"user":     types.StringValue(service.Admin.User),
			"password": types.StringValue(service.Admin.Password),
		},
		diags,
	)
	data.DatabaseAdmin = utils.ObjectValue(
		map[string]attr.Type{
			"host":     types.StringType,
			"port":     types.StringType,
			"user":     types.StringType,
			"password": types.StringType,
			"command":  types.StringType,
		},
		map[string]attr.Value{
			"host":     types.StringValue(service.DatabaseAdmin.Host),
			"port":     types.StringValue(service.DatabaseAdmin.Port),
			"user":     types.StringValue(service.DatabaseAdmin.User),
			"password": types.StringValue(service.DatabaseAdmin.Password),
			"command":  types.StringValue(service.DatabaseAdmin.Command),
		},
		diags,
	)
	data.GlobalIP = types.StringValue(service.GlobalIP)
	data.TrafficOutgoing = types.Int64Value(service.TrafficOutgoing)
	data.TrafficIncoming = types.Int64Value(service.TrafficIncoming)
	data.TrafficIncluded = types.Int64Value(service.TrafficIncluded)
	data.Cores = types.Int64Value(service.Cores)
	data.RAMSizeGB = types.StringValue(service.RAMSizeGB)
	data.StorageSizeGB = types.Int64Value(service.StorageSizeGB)
	data.PricePerHour = types.StringValue(service.PricePerHour)
	data.AppAutoUpdatesEnabled = utils.BoolValue(service.AppAutoUpdatesEnabled)
	data.AppAutoUpdatesDayOfWeek = types.Int64Value(service.AppAutoUpdatesDayOfWeek)
	data.AppAutoUpdatesHour = types.Int64Value(service.AppAutoUpdatesHour)
	data.AppAutoUpdatesMinute = types.Int64Value(service.AppAutoUpdatesMinute)
	data.SystemAutoUpdatesEnabled = utils.BoolValue(service.SystemAutoUpdatesEnabled)
	if !data.SystemAutoUpdatesEnabled.ValueBool() {
		// If system auto updates are disabled, then security patches only is also disabled
		data.SystemAutoUpdatesSecurityPatchesOnlyEnabled = types.BoolValue(false)
	} else {
		data.SystemAutoUpdatesSecurityPatchesOnlyEnabled = utils.BoolValue(service.SystemAutoUpdatesSecurityPatchesOnlyEnabled)
	}
	data.SystemAutoUpdatesRebootDayOfWeek = types.Int64Value(service.SystemAutoUpdatesRebootDayOfWeek)
	data.SystemAutoUpdatesRebootHour = types.Int64Value(service.SystemAutoUpdatesRebootHour)
	data.SystemAutoUpdatesRebootMinute = types.Int64Value(service.SystemAutoUpdatesRebootMinute)
	data.BackupsEnabled = utils.BoolValue(service.BackupsEnabled)
	data.RemoteBackupsEnabled = utils.BoolValue(service.RemoteBackupsEnabled)
	data.ExternalBackupsEnabled = utils.BoolValue(service.ExternalBackupsEnabled)
	data.ExternalBackupsUpdateDayOfWeek = types.Int64Value(service.ExternalBackupsUpdateDayOfWeek)
	data.ExternalBackupsUpdateHour = types.Int64Value(service.ExternalBackupsUpdateHour)
	data.ExternalBackupsUpdateMinute = types.Int64Value(service.ExternalBackupsUpdateMinute)
	data.ExternalBackupsUpdateType = types.StringValue(service.ExternalBackupsUpdateType)
	data.ExternalBackupsRetainDayOfWeek = types.Int64Value(service.ExternalBackupsRetainDayOfWeek)
	data.FirewallEnabled = utils.BoolValue(service.FirewallEnabled)
	data.FirewallId = types.StringValue(service.FirewallID)
	data.FirewallPorts = types.StringValue(service.FirewallPorts)
	data.AlertsEnabled = utils.BoolValue(service.AlertsEnabled)
	data.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
}

func (r *ServiceResource) waitServiceCreate(ctx context.Context, service *elestio.Service) (*elestio.Service, error) {
	createTimeout := 20 * time.Minute
	createStateConf := sdk_resource.StateChangeConf{
		Pending: []string{"creating"},
		Target:  []string{"created"},
		Refresh: func() (interface{}, string, error) {
			serviceW, err := r.client.Service.Get(service.ProjectID, service.ID)
			if err != nil {
				return struct{}{}, "", err
			}

			if serviceW.DeploymentStatus != elestio.ServiceDeploymentStatusDeployed {
				return struct{}{}, "creating", nil
			}

			// App auto updates are enabled by default at service creation
			if !utils.BoolValue(serviceW.AppAutoUpdatesEnabled).ValueBool() {
				return struct{}{}, "creating", nil
			}

			// System auto updates are enabled by default at service creation
			if !utils.BoolValue(serviceW.SystemAutoUpdatesEnabled).ValueBool() {
				return struct{}{}, "creating", nil
			}

			// Backups are enabled by default at service creation if service level is greater than level1
			if serviceW.SupportLevel != "level1" && !utils.BoolValue(serviceW.BackupsEnabled).ValueBool() {
				return struct{}{}, "creating", nil
			}

			// Remote backups are enabled by default at service creation
			if !utils.BoolValue(serviceW.RemoteBackupsEnabled).ValueBool() {
				return struct{}{}, "creating", nil
			}

			// Firewall is enabled by default at service creation
			if !utils.BoolValue(serviceW.FirewallEnabled).ValueBool() {
				return struct{}{}, "creating", nil
			}

			// Alerts are enabled by default at service creation
			if !utils.BoolValue(serviceW.AlertsEnabled).ValueBool() {
				return struct{}{}, "creating", nil
			}

			return serviceW, "created", nil
		},
		Timeout:                   createTimeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	tflog.Trace(ctx, fmt.Sprintf("Service creation waiter timeout %.0f minutes", createTimeout.Minutes()))

	serviceCreated, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("service creation waiter failed, got error: %s", err)
	}

	return serviceCreated.(*elestio.Service), nil
}

func (r *ServiceResource) waitServiceDelete(ctx context.Context, service *elestio.Service) error {
	deleteTimeout := 20 * time.Minute
	deleteStateConf := sdk_resource.StateChangeConf{
		Pending: []string{"deleting"},
		Target:  []string{"deleted"},
		Refresh: func() (any, string, error) {
			_, err := r.client.Service.Get(service.ProjectID, service.ID)

			if err == nil {
				return struct{}{}, "deleting", nil
			}

			return struct{}{}, "deleted", nil
		},
		Timeout:                   deleteTimeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	tflog.Trace(ctx, fmt.Sprintf("Service deletion waiter timeout %.0f minutes", deleteTimeout.Minutes()))

	if _, err := deleteStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("service deletion waiter failed, got error: %s", err)
	}

	return nil
}

func (r *ServiceResource) waitServerTypeUpdate(ctx context.Context, service *elestio.Service, expectedNewServerType string) (*elestio.Service, error) {
	timeout := 10 * time.Minute
	createStateConf := sdk_resource.StateChangeConf{
		Pending: []string{"updating"},
		Target:  []string{"updated"},
		Refresh: func() (interface{}, string, error) {
			serviceW, err := r.client.Service.Get(service.ProjectID, service.ID)
			if err != nil {
				return struct{}{}, "", err
			}

			// running -> stopped -> migrating -> running
			if serviceW.Status != elestio.ServiceStatusRunning {
				return struct{}{}, "updating", nil
			}

			if serviceW.ServerType != expectedNewServerType {
				return struct{}{}, "updating", nil
			}

			return serviceW, "updated", nil
		},
		Timeout:                   timeout,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	tflog.Trace(ctx, fmt.Sprintf("Service creation waiter timeout %.0f minutes", timeout.Minutes()))

	serviceCreated, err := createStateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("service creation waiter failed, got error: %s", err)
	}

	return serviceCreated.(*elestio.Service), nil
}
