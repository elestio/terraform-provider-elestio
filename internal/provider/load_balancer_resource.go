package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elestio/elestio-go-api-client"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
	Id           types.String `tfsdk:"id"`
	ProjectId    types.String `tfsdk:"project_id"`
	ProviderName types.String `tfsdk:"provider_name"`
	Datacenter   types.String `tfsdk:"datacenter"`
	ServerType   types.String `tfsdk:"server_type"`
	Config       ConfigModel  `tfsdk:"config"`
}

type (
	ConfigModel struct {
		HostHeader             types.String       `tfsdk:"host_header"`
		IsAccessLogsEnabled    types.Bool         `tfsdk:"is_access_logs_enabled"`
		IsForceHTTPSEnabled    types.Bool         `tfsdk:"is_force_https_enabled"`
		IPRateLimit            types.Int64        `tfsdk:"ip_rate_limit"`
		IsIPRateLimitEnabled   types.Bool         `tfsdk:"is_ip_rate_limit_enabled"`
		OutputCacheInSeconds   types.Int64        `tfsdk:"output_cache_in_seconds"`
		IsStickySessionEnabled types.Bool         `tfsdk:"is_sticky_session_enabled"`
		IsProxyProtocolEnabled types.Bool         `tfsdk:"is_proxy_protocol_enabled"`
		SSLDomains             types.Set          `tfsdk:"ssl_domains"`
		ForwardRules           []ForwardRuleModel `tfsdk:"forward_rules"`
		// OutputHeaders          []struct {
		// 	Key   types.String `tfsdk:"key"`
		// 	Value types.String `tfsdk:"value"`
		// } `tfsdk:"output_headers"`
		TargetServices        types.Set `tfsdk:"target_services"`
		RemoveResponseHeaders types.Set `tfsdk:"remove_response_headers"`
	}

	ForwardRuleModel struct {
		Port           types.String `tfsdk:"port"`
		Protocol       types.String `tfsdk:"protocol"`
		TargetPort     types.String `tfsdk:"target_port"`
		TargetProtocol types.String `tfsdk:"target_protocol"`
	}
)

var (
	configAttrTypes = map[string]attr.Type{
		"host_header":               types.StringType,
		"is_access_logs_enabled":    types.BoolType,
		"is_force_https_enabled":    types.BoolType,
		"ip_rate_limit":             types.Int64Type,
		"is_ip_rate_limit_enabled":  types.BoolType,
		"output_cache_in_seconds":   types.Int64Type,
		"is_sticky_session_enabled": types.BoolType,
		"is_proxy_protocol_enabled": types.BoolType,
		"ssl_domains":               types.SetType{ElemType: types.StringType},
		"forward_rules":             types.SetType{ElemType: types.ObjectType{AttrTypes: forwardRuleAttrTypes}},
	}

	forwardRuleAttrTypes = map[string]attr.Type{
		"port":            types.StringType,
		"protocol":        types.StringType,
		"target_port":     types.StringType,
		"target_protocol": types.StringType,
	}
)

func (r *LoadBalancerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer"
}

func (r *LoadBalancerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	var diags diag.Diagnostics

	defaultSSLDomains, diags := types.SetValue(types.StringType, []attr.Value{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultForwardRules, diags := types.SetValue(types.ObjectType{AttrTypes: forwardRuleAttrTypes}, []attr.Value{
		types.ObjectValueMust(forwardRuleAttrTypes, map[string]attr.Value{
			"port":            types.StringValue("80"),
			"protocol":        types.StringValue("HTTP"),
			"target_port":     types.StringValue("3000"),
			"target_protocol": types.StringValue("HTTP"),
		}),
		types.ObjectValueMust(forwardRuleAttrTypes, map[string]attr.Value{
			"port":            types.StringValue("443"),
			"protocol":        types.StringValue("HTTPS"),
			"target_port":     types.StringValue("3000"),
			"target_protocol": types.StringValue("HTTP"),
		}),
	})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultRemoveResponseHeaders, diags := types.SetValue(types.StringType, []attr.Value{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Provider name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"datacenter": schema.StringAttribute{
				MarkdownDescription: "Datacenter name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Server type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.SingleNestedAttribute{
				MarkdownDescription: "Load balancer configuration",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"host_header": schema.StringAttribute{
						MarkdownDescription: "Host header",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("$http_host"),
					},
					"is_access_logs_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is access logs enabled",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"is_force_https_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is force https enabled",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
					"ip_rate_limit": schema.Int64Attribute{
						MarkdownDescription: "IP rate limit (requests per second)",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(100),
					},
					"is_ip_rate_limit_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is IP rate limit enabled",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"output_cache_in_seconds": schema.Int64Attribute{
						MarkdownDescription: "Output cache in seconds",
						Optional:            true,
						Computed:            true,
						Default:             int64default.StaticInt64(0),
					},
					"is_sticky_session_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is sticky session enabled",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"is_proxy_protocol_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is proxy protocol enabled",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(false),
					},
					"ssl_domains": schema.SetAttribute{
						MarkdownDescription: "SSL domains",
						Optional:            true,
						Computed:            true,
						Default:             setdefault.StaticValue(defaultSSLDomains),
						ElementType:         types.StringType,
					},
					"forward_rules": schema.SetNestedAttribute{
						MarkdownDescription: "Forward rules",
						Optional:            true,
						Computed:            true,
						Default:             setdefault.StaticValue(defaultForwardRules),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"protocol": schema.StringAttribute{
									MarkdownDescription: "Protocol",
									Required:            true,
								},
								"port": schema.StringAttribute{
									MarkdownDescription: "Port",
									Required:            true,
								},
								"target_protocol": schema.StringAttribute{
									MarkdownDescription: "Target protocol",
									Required:            true,
								},
								"target_port": schema.StringAttribute{
									MarkdownDescription: "Target port",
									Required:            true,
								},
							},
						},
					},
					// "output_headers": schema.ListNestedAttribute{
					// 	MarkdownDescription: "Output headers",
					// 	Optional:            true,
					// 	Computed:            true,
					// 	NestedObject: schema.NestedAttributeObject{
					// 		Attributes: map[string]schema.Attribute{
					// 			"key": schema.StringAttribute{
					// 				MarkdownDescription: "Key",
					// 				Required:            true,
					// 			},
					// 			"value": schema.StringAttribute{
					// 				MarkdownDescription: "Value",
					// 				Required:            true,
					// 			},
					// 		},
					// 	},
					// },
					"target_services": schema.SetAttribute{
						MarkdownDescription: "Target services",
						Required:            true,
						ElementType:         types.StringType,
						Validators: []validator.Set{
							setvalidator.SizeAtLeast(1),
						},
					},
					"remove_response_headers": schema.SetAttribute{
						MarkdownDescription: "Remove response headers",
						Optional:            true,
						Computed:            true,
						Default:             setdefault.StaticValue(defaultRemoveResponseHeaders),
						ElementType:         types.StringType,
					},
				},
			},
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

	var diags diag.Diagnostics

	sslDomains := []string{}
	diags = plan.Config.SSLDomains.ElementsAs(ctx, &sslDomains, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	forwardRules := make([]elestio.LoadBalancerConfigForwardRule, len(plan.Config.ForwardRules))
	for i, forwardRule := range plan.Config.ForwardRules {
		forwardRules[i] = elestio.LoadBalancerConfigForwardRule{
			Protocol:       forwardRule.Protocol.ValueString(),
			Port:           forwardRule.Port.ValueString(),
			TargetProtocol: forwardRule.TargetProtocol.ValueString(),
			TargetPort:     forwardRule.TargetPort.ValueString(),
		}
	}

	targetServices := []string{}
	diags = plan.Config.TargetServices.ElementsAs(ctx, &targetServices, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	removeResponseHeaders := []string{}
	diags = plan.Config.RemoveResponseHeaders.ElementsAs(ctx, &removeResponseHeaders, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	clientLoadBalancer, err := r.client.LoadBalancer.Create(elestio.CreateLoadBalancerRequest{
		ProjectID:    plan.ProjectId.ValueString(),
		ProviderName: plan.ProviderName.ValueString(),
		Datacenter:   plan.Datacenter.ValueString(),
		ServerType:   plan.ServerType.ValueString(),
		Config: elestio.CreateLoadBalancerRequestConfig{
			HostHeader:             plan.Config.HostHeader.ValueString(),
			IsAccessLogsEnabled:    plan.Config.IsAccessLogsEnabled.ValueBool(),
			IsForceHTTPSEnabled:    plan.Config.IsForceHTTPSEnabled.ValueBool(),
			IPRateLimit:            plan.Config.IPRateLimit.ValueInt64(),
			IsIPRateLimitEnabled:   plan.Config.IsIPRateLimitEnabled.ValueBool(),
			OutputCacheInSeconds:   plan.Config.OutputCacheInSeconds.ValueInt64(),
			IsStickySessionEnabled: plan.Config.IsStickySessionEnabled.ValueBool(),
			IsProxyProtocolEnabled: plan.Config.IsProxyProtocolEnabled.ValueBool(),
			SSLDomains:             sslDomains,
			ForwardRules:           forwardRules,
			OutputHeaders:          []elestio.LoadBalancerConfigOutputHeader{},
			TargetServices:         targetServices,
			RemoveResponseHeaders:  removeResponseHeaders,
		},
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create load balancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created Load Balancer: "+clientLoadBalancer.ID)

	loadBalancer, d := transformClientLoadBalancerToResourceModel(ctx, clientLoadBalancer, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &loadBalancer)...)
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
	loadBalancerModel, d := transformClientLoadBalancerToResourceModel(ctx, loadBalancer, state)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	loadBalancerModel, d := transformClientLoadBalancerToResourceModel(ctx, loadBalancer, plan)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
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

	confirmDeleteStateConf := &retry.StateChangeConf{
		Pending: []string{"DELETING"},
		Target:  []string{"DELETED"},
		Refresh: func() (interface{}, string, error) {
			loadBalancer, err := r.client.LoadBalancer.Get(state.ProjectId.ValueString(), state.Id.ValueString())
			// We expect a 401 error when the load balancer is deleted
			if err != nil {
				return struct{}{}, "DELETED", nil
			}
			return loadBalancer, "DELETING", nil
		},
		Timeout:                   5 * time.Minute,
		Delay:                     10 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	_, err = confirmDeleteStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to confirm load balancer deletion, got error: %s", err))
	}
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

func transformClientLoadBalancerToResourceModel(ctx context.Context, loadBalancerClient *elestio.LoadBalancer, state *LoadBalancerResourceModel) (LoadBalancerResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	var loadBalancerModel LoadBalancerResourceModel
	loadBalancerModel.Id = types.StringValue(loadBalancerClient.ID)
	loadBalancerModel.ProjectId = types.StringValue(loadBalancerClient.ProjectID)
	loadBalancerModel.ProviderName = types.StringValue(loadBalancerClient.ProviderName)
	loadBalancerModel.Datacenter = types.StringValue(loadBalancerClient.Datacenter)
	loadBalancerModel.ServerType = types.StringValue(loadBalancerClient.ServerType)
	loadBalancerModel.Config.HostHeader = types.StringValue(loadBalancerClient.Config.HostHeader)
	loadBalancerModel.Config.IsAccessLogsEnabled = types.BoolValue(loadBalancerClient.Config.IsAccessLogsEnabled)
	loadBalancerModel.Config.IsForceHTTPSEnabled = types.BoolValue(loadBalancerClient.Config.IsForceHTTPSEnabled)
	loadBalancerModel.Config.IPRateLimit = types.Int64Value(loadBalancerClient.Config.IPRateLimit)
	loadBalancerModel.Config.IsIPRateLimitEnabled = types.BoolValue(loadBalancerClient.Config.IsIPRateLimitEnabled)
	loadBalancerModel.Config.OutputCacheInSeconds = types.Int64Value(loadBalancerClient.Config.OutputCacheInSeconds)
	loadBalancerModel.Config.IsStickySessionEnabled = types.BoolValue(loadBalancerClient.Config.IsStickySessionEnabled)
	loadBalancerModel.Config.IsProxyProtocolEnabled = types.BoolValue(loadBalancerClient.Config.IsProxyProtocolEnabled)

	loadBalancerModel.Config.SSLDomains, diags = types.SetValueFrom(ctx, types.StringType, loadBalancerClient.Config.SSLDomains)
	if diags.HasError() {
		return loadBalancerModel, diags
	}

	var forwardRules []ForwardRuleModel
	for _, r := range loadBalancerClient.Config.ForwardRules {
		var rule ForwardRuleModel
		rule.Port = types.StringValue(r.Port)
		rule.Protocol = types.StringValue(r.Protocol)
		rule.TargetPort = types.StringValue(r.TargetPort)
		rule.TargetProtocol = types.StringValue(r.TargetProtocol)
		forwardRules = append(forwardRules, rule)
	}
	loadBalancerModel.Config.ForwardRules = forwardRules

	loadBalancerModel.Config.TargetServices, diags = types.SetValueFrom(ctx, types.StringType, loadBalancerClient.Config.TargetServices)
	if diags.HasError() {
		return loadBalancerModel, diags
	}

	loadBalancerModel.Config.RemoveResponseHeaders, diags = types.SetValueFrom(ctx, types.StringType, loadBalancerClient.Config.RemoveResponseHeaders)
	if diags.HasError() {
		return loadBalancerModel, diags
	}
	return loadBalancerModel, diags
}

// func (m LoadBalancerResourceModel) ToClientType(ctx context.Context) (elestio.LoadBalancer, diag.Diagnostics) {
// 	var diags diag.Diagnostics
// 	var loadBalancerClient elestio.LoadBalancer

// 	loadBalancerClient.ID = m.Id.ValueString()
// 	loadBalancerClient.ProjectID = m.ProjectId.ValueString()
// 	loadBalancerClient.ProviderName = m.ProviderName.ValueString()
// 	loadBalancerClient.Datacenter = m.Datacenter.ValueString()
// 	loadBalancerClient.ServerType = m.ServerType.ValueString()
// 	loadBalancerClient.Config.HostHeader = m.Config.HostHeader.ValueString()
// 	loadBalancerClient.Config.IsAccessLogsEnabled = m.Config.IsAccessLogsEnabled.ValueBool()
// 	loadBalancerClient.Config.IsForceHTTPSEnabled = m.Config.IsForceHTTPSEnabled.ValueBool()
// 	loadBalancerClient.Config.IPRateLimit = m.Config.IPRateLimit.ValueInt64()
// 	loadBalancerClient.Config.IsIPRateLimitEnabled = m.Config.IsIPRateLimitEnabled.ValueBool()
// 	loadBalancerClient.Config.OutputCacheInSeconds = m.Config.OutputCacheInSeconds.ValueInt64()
// 	loadBalancerClient.Config.IsStickySessionEnabled = m.Config.IsStickySessionEnabled.ValueBool()
// 	loadBalancerClient.Config.IsProxyProtocolEnabled = m.Config.IsProxyProtocolEnabled.ValueBool()
// 	// config.SSLDomains = utils.SetValueOrNull(ctx, types.StringType, l.Config.SSLDomains, &diags)
// 	var forwardRules []elestio.LoadBalancerConfigForwardRule
// 	for _, rule := range m.Config.ForwardRules {
// 		var forwardRule elestio.LoadBalancerConfigForwardRule
// 		forwardRule.Port = rule.Port.ValueString()
// 		forwardRule.Protocol = rule.Protocol.ValueString()
// 		forwardRule.TargetPort = rule.TargetPort.ValueString()
// 		forwardRule.TargetProtocol = rule.TargetProtocol.ValueString()
// 		forwardRules = append(forwardRules, forwardRule)
// 	}
// 	loadBalancerClient.Config.ForwardRules = forwardRules

// 	return loadBalancerClient, diags
// }

// func (m ConfigModel) ToClientTypeConfig(ctx context.Context) (elestio.LoadBalancerConfig, diag.Diagnostics) {
// 	var diags diag.Diagnostics
// 	var result elestio.LoadBalancerConfig

// 	result.HostHeader = m.HostHeader.ValueString()
// 	result.IsAccessLogsEnabled = m.IsAccessLogsEnabled.ValueBool()
// 	result.IsForceHTTPSEnabled = m.IsForceHTTPSEnabled.ValueBool()
// 	result.IPRateLimit = m.IPRateLimit.ValueInt64()
// 	result.IsIPRateLimitEnabled = m.IsIPRateLimitEnabled.ValueBool()
// 	result.OutputCacheInSeconds = m.OutputCacheInSeconds.ValueInt64()
// 	result.IsStickySessionEnabled = m.IsStickySessionEnabled.ValueBool()
// 	result.IsProxyProtocolEnabled = m.IsProxyProtocolEnabled.ValueBool()

// 	var forwardRulesObjects []types.Object
// 	diags.Append(m.ForwardRules.ElementsAs(ctx, &forwardRulesObjects, false)...)
// 	if diags.HasError() {
// 		return result, diags
// 	}

// 	for _, forwardRuleObject := range forwardRulesObjects {
// 		var forwardRuleModel ForwardRuleModel
// 		diags.Append(forwardRuleObject.As(ctx, &forwardRuleModel, basetypes.ObjectAsOptions{})...)
// 		if diags.HasError() {
// 			return result, diags
// 		}

// 		forwardRule, forwardRuleDiags := forwardRuleModel.ToClientTypeForwardRule(ctx)
// 		diags.Append(forwardRuleDiags...)
// 		if diags.HasError() {
// 			return result, diags
// 		}

// 		result.ForwardRules = append(result.ForwardRules, forwardRule)
// 	}

// 	return result, diags
// }

// func (m ForwardRuleModel) ToClientTypeForwardRule(ctx context.Context) (elestio.LoadBalancerConfigForwardRule, diag.Diagnostics) {
// 	var diags diag.Diagnostics
// 	var result elestio.LoadBalancerConfigForwardRule

// 	result.Port = m.Port.ValueString()
// 	result.Protocol = m.Protocol.ValueString()
// 	result.TargetPort = m.TargetPort.ValueString()
// 	result.TargetProtocol = m.TargetProtocol.ValueString()

// 	return result, diags
// }
