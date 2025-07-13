package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elestio/elestio-go-api-client/v2"
	"github.com/elestio/terraform-provider-elestio/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	Id               types.String `tfsdk:"id"`
	ProjectId        types.String `tfsdk:"project_id"`
	ProviderName     types.String `tfsdk:"provider_name"`
	Datacenter       types.String `tfsdk:"datacenter"`
	ServerType       types.String `tfsdk:"server_type"`
	Config           ConfigModel  `tfsdk:"config"`
	CreatedAt        types.String `tfsdk:"created_at"`
	CreatedBy        types.String `tfsdk:"created_by"`
	DeploymentStatus types.String `tfsdk:"deployment_status"`
	IPV4             types.String `tfsdk:"ipv4"`
	IPV6             types.String `tfsdk:"ipv6"`
	CNAME            types.String `tfsdk:"cname"`
	Country          types.String `tfsdk:"country"`
	City             types.String `tfsdk:"city"`
	GlobalIP         types.String `tfsdk:"global_ip"`
	Cores            types.Int64  `tfsdk:"cores"`
	RAMSizeGB        types.String `tfsdk:"ram_size_gb"`
	StorageSizeGB    types.Int64  `tfsdk:"storage_size_gb"`
	PricePerHour     types.String `tfsdk:"price_per_hour"`
}

type ConfigModel struct {
	HostHeader            types.String        `tfsdk:"host_header"`
	AccessLogsEnabled     types.Bool          `tfsdk:"access_logs_enabled"`
	ForceHTTPSEnabled     types.Bool          `tfsdk:"force_https_enabled"`
	IPRateLimitPerSecond  types.Int64         `tfsdk:"ip_rate_limit_per_second"`
	IPRateLimitEnabled    types.Bool          `tfsdk:"ip_rate_limit_enabled"`
	OutputCacheInSeconds  types.Int64         `tfsdk:"output_cache_in_seconds"`
	StickySessionEnabled  types.Bool          `tfsdk:"sticky_session_enabled"`
	ProxyProtocolEnabled  types.Bool          `tfsdk:"proxy_protocol_enabled"`
	SSLDomains            types.Set           `tfsdk:"ssl_domains"`
	ForwardRules          []ForwardRuleModel  `tfsdk:"forward_rules"`
	OutputHeaders         []OutputHeaderModel `tfsdk:"output_headers"`
	TargetServices        types.Set           `tfsdk:"target_services"`
	RemoveResponseHeaders types.Set           `tfsdk:"remove_response_headers"`
}

// var configAttrTypes = map[string]attr.Type{
// 	"host_header":              types.StringType,
// 	"access_logs_enabled":      types.BoolType,
// 	"force_https_enabled":      types.BoolType,
// 	"ip_rate_limit_per_second": types.Int64Type,
// 	"ip_rate_limit_enabled":    types.BoolType,
// 	"output_cache_in_seconds":  types.Int64Type,
// 	"sticky_session_enabled":   types.BoolType,
// 	"proxy_protocol_enabled":   types.BoolType,
// 	"ssl_domains":              types.SetType{ElemType: types.StringType},
// 	"forward_rules":            types.SetType{ElemType: types.ObjectType{AttrTypes: forwardRuleAttrTypes}},
// 	"output_headers":           types.SetType{ElemType: types.ObjectType{AttrTypes: outputHeaderAttrTypes}},
// 	"target_services":          types.SetType{ElemType: types.StringType},
// 	"remove_response_headers":  types.SetType{ElemType: types.StringType},
// }

type ForwardRuleModel struct {
	Port           types.String `tfsdk:"port"`
	Protocol       types.String `tfsdk:"protocol"`
	TargetPort     types.String `tfsdk:"target_port"`
	TargetProtocol types.String `tfsdk:"target_protocol"`
}

var forwardRuleAttrTypes = map[string]attr.Type{
	"port":            types.StringType,
	"protocol":        types.StringType,
	"target_port":     types.StringType,
	"target_protocol": types.StringType,
}

type OutputHeaderModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

var outputHeaderAttrTypes = map[string]attr.Type{
	"key":   types.StringType,
	"value": types.StringType,
}

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

	defaultOutputHeaders, diags := types.SetValue(types.ObjectType{AttrTypes: outputHeaderAttrTypes}, []attr.Value{})
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
				MarkdownDescription: "Provider name." +
					" Availables values on the related guide: https://docs.elest.io/books/elestio-terraform-provider/page/providers-datacenters-and-server-types",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("hetzner", "do", "lightsail", "linode", "vultr", "scaleway"),
				},
			},
			"datacenter": schema.StringAttribute{
				MarkdownDescription: "Datacenter name." +
					" Availables values on the related guide: https://docs.elest.io/books/elestio-terraform-provider/page/providers-datacenters-and-server-types",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Server type." +
					" Availables values on the related guide: https://docs.elest.io/books/elestio-terraform-provider/page/providers-datacenters-and-server-types",
				Required: true,
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
						MarkdownDescription: "Host header." +
							"</br>Default value: `$http_host`",
						Optional: true,
						Computed: true,
						Default:  stringdefault.StaticString("$http_host"),
					},
					"access_logs_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is access logs enabled." +
							"</br>Default value: `true`",
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"force_https_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is force https enabled." +
							" </br>Default value: `true`",
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(true),
					},
					"ip_rate_limit_per_second": schema.Int64Attribute{
						MarkdownDescription: "Indicate the maximum number of requests allowed per second per IP address." +
							" </br>Default value: `100`",
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(100),
					},
					"ip_rate_limit_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is IP rate limit enabled." +
							" </br>Default value: `false`",
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"output_cache_in_seconds": schema.Int64Attribute{
						MarkdownDescription: "Output cache in seconds." +
							" </br>Default value: `0`",
						Optional: true,
						Computed: true,
						Default:  int64default.StaticInt64(0),
					},
					"sticky_session_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is sticky session enabled." +
							" </br>Default value: `false`",
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"proxy_protocol_enabled": schema.BoolAttribute{
						MarkdownDescription: "Is proxy protocol enabled." +
							" </br>Default value: `false`",
						Optional: true,
						Computed: true,
						Default:  booldefault.StaticBool(false),
					},
					"ssl_domains": schema.SetAttribute{
						MarkdownDescription: "SSL domains",
						Optional:            true,
						Computed:            true,
						Default:             setdefault.StaticValue(defaultSSLDomains),
						ElementType:         types.StringType,
					},
					"forward_rules": schema.SetNestedAttribute{
						MarkdownDescription: "Forward Rules." +
							fmt.Sprintf(" </br>Default value: `%+v`", defaultForwardRules),
						Optional: true,
						Computed: true,
						Default:  setdefault.StaticValue(defaultForwardRules),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"protocol": schema.StringAttribute{
									MarkdownDescription: "Protocol. Availables values: `HTTP`, `HTTPS`, `TCP`, `UDP`. If you use `TCP` or `UDP`, you must use the same value for relative `target_protocol` attribute.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("HTTP", "HTTPS", "TCP", "UDP"),
										validators.IsValidForwardRuleMatchingProtocol(path.MatchRelative().AtParent().AtName("target_protocol"), "TCP", "UDP"),
									},
								},
								"port": schema.StringAttribute{
									MarkdownDescription: "Port",
									Required:            true,
								},
								"target_protocol": schema.StringAttribute{
									MarkdownDescription: "Target protocol. Availables values: `HTTP`, `HTTPS`, `TCP`, `UDP`. If you use `TCP` or `UDP`, you must use the same value for relative `protocol` attribute.",
									Required:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("HTTP", "HTTPS", "TCP", "UDP"),
										validators.IsValidForwardRuleMatchingProtocol(path.MatchRelative().AtParent().AtName("protocol"), "TCP", "UDP"),
									},
								},
								"target_port": schema.StringAttribute{
									MarkdownDescription: "Target port",
									Required:            true,
								},
							},
						},
					},
					"output_headers": schema.SetNestedAttribute{
						MarkdownDescription: "Output headers",
						Optional:            true,
						Computed:            true,
						Default:             setdefault.StaticValue(defaultOutputHeaders),
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"key": schema.StringAttribute{
									MarkdownDescription: "Key",
									Required:            true,
								},
								"value": schema.StringAttribute{
									MarkdownDescription: "Value",
									Required:            true,
								},
							},
						},
					},
					"target_services": schema.SetAttribute{
						MarkdownDescription: " The services to which the load balancer will forward the requests." +
							" You can provide services IDs but also IPs and CNAME records." +
							"</br>Example: `[\"xxxx-xxxx-xxxx-xxxx\", \"192.168.xxxx\", \"myawesomeapp.com\"]`",
						Required:    true,
						ElementType: types.StringType,
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
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Created at",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created_by": schema.StringAttribute{
				MarkdownDescription: "Created by",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"deployment_status": schema.StringAttribute{
				MarkdownDescription: "Deployment status",
				Computed:            true,
			},
			"ipv4": schema.StringAttribute{
				MarkdownDescription: "IPv4",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv6": schema.StringAttribute{
				MarkdownDescription: "IPv6",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cname": schema.StringAttribute{
				MarkdownDescription: "CNAME",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "Country",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "City",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"global_ip": schema.StringAttribute{
				MarkdownDescription: "Global IP",
				Computed:            true,
			},
			"cores": schema.Int64Attribute{
				MarkdownDescription: "Cores",
				Computed:            true,
			},
			"ram_size_gb": schema.StringAttribute{
				MarkdownDescription: "RAM size in GB",
				Computed:            true,
			},
			"storage_size_gb": schema.Int64Attribute{
				MarkdownDescription: "Storage size in GB",
				Computed:            true,
			},
			"price_per_hour": schema.StringAttribute{
				MarkdownDescription: "Price per hour",
				Computed:            true,
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

	planSSLDomains := []string{}
	diags = plan.Config.SSLDomains.ElementsAs(ctx, &planSSLDomains, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	planForwardRules := make([]elestio.LoadBalancerConfigForwardRule, len(plan.Config.ForwardRules))
	for i, forwardRule := range plan.Config.ForwardRules {
		planForwardRules[i] = elestio.LoadBalancerConfigForwardRule{
			Protocol:       forwardRule.Protocol.ValueString(),
			Port:           forwardRule.Port.ValueString(),
			TargetProtocol: forwardRule.TargetProtocol.ValueString(),
			TargetPort:     forwardRule.TargetPort.ValueString(),
		}
	}

	planOutputHeaders := make([]elestio.LoadBalancerConfigOutputHeader, len(plan.Config.OutputHeaders))
	for i, outputHeader := range plan.Config.OutputHeaders {
		planOutputHeaders[i] = elestio.LoadBalancerConfigOutputHeader{
			Key:   outputHeader.Key.ValueString(),
			Value: outputHeader.Value.ValueString(),
		}
	}

	planTargetServices := []string{}
	diags = plan.Config.TargetServices.ElementsAs(ctx, &planTargetServices, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	planRemoveResponseHeaders := []string{}
	diags = plan.Config.RemoveResponseHeaders.ElementsAs(ctx, &planRemoveResponseHeaders, false)
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
			IsAccessLogsEnabled:    plan.Config.AccessLogsEnabled.ValueBool(),
			IsForceHTTPSEnabled:    plan.Config.ForceHTTPSEnabled.ValueBool(),
			IPRateLimit:            plan.Config.IPRateLimitPerSecond.ValueInt64(),
			IsIPRateLimitEnabled:   plan.Config.IPRateLimitEnabled.ValueBool(),
			OutputCacheInSeconds:   plan.Config.OutputCacheInSeconds.ValueInt64(),
			IsStickySessionEnabled: plan.Config.StickySessionEnabled.ValueBool(),
			IsProxyProtocolEnabled: plan.Config.ProxyProtocolEnabled.ValueBool(),
			SSLDomains:             planSSLDomains,
			ForwardRules:           planForwardRules,
			OutputHeaders:          planOutputHeaders,
			TargetServices:         planTargetServices,
			RemoveResponseHeaders:  planRemoveResponseHeaders,
		},
		CreatedFrom: "terraform",
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create load balancer, got error: %s", err))
		return
	}

	tflog.Trace(ctx, "Created Load Balancer: "+clientLoadBalancer.ID)

	tflog.Info(ctx, "Waiting for Load Balancer to be deployed: "+clientLoadBalancer.ID)
	deploymentStateConf := &retry.StateChangeConf{
		Pending: []string{"DEPLOYING"},
		Target:  []string{"DEPLOYED"},
		Refresh: func() (interface{}, string, error) {
			deployingClientLoadBalancer, err := r.client.LoadBalancer.Get(clientLoadBalancer.ProjectID, clientLoadBalancer.ID)
			if err != nil {
				return nil, "", fmt.Errorf("waiting for load balancer deployment, got error: %s", err)
			}

			if deployingClientLoadBalancer.DeploymentStatus != elestio.LoadBalancerDeploymentStatusDeployed {
				return deployingClientLoadBalancer, "DEPLOYING", nil
			}

			return deployingClientLoadBalancer, "DEPLOYED", nil
		},
		Timeout:                   10 * time.Minute,
		Delay:                     30 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}
	deployedClientLoadBalancer, err := deploymentStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to wait for load balancer deployment, got error: %s", err))
		return
	}
	tflog.Trace(ctx, "Load Balancer deployed: "+deployedClientLoadBalancer.(*elestio.LoadBalancer).ID)

	loadBalancer, d := transformClientLoadBalancerToResourceModel(ctx, deployedClientLoadBalancer.(*elestio.LoadBalancer), plan)
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

	var diags diag.Diagnostics

	planSSLDomains := []string{}
	diags = plan.Config.SSLDomains.ElementsAs(ctx, &planSSLDomains, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	planForwardRules := make([]elestio.LoadBalancerConfigForwardRule, len(plan.Config.ForwardRules))
	for i, forwardRule := range plan.Config.ForwardRules {
		planForwardRules[i] = elestio.LoadBalancerConfigForwardRule{
			Protocol:       forwardRule.Protocol.ValueString(),
			Port:           forwardRule.Port.ValueString(),
			TargetProtocol: forwardRule.TargetProtocol.ValueString(),
			TargetPort:     forwardRule.TargetPort.ValueString(),
		}
	}

	planOutputHeaders := make([]elestio.LoadBalancerConfigOutputHeader, len(plan.Config.OutputHeaders))
	for i, outputHeader := range plan.Config.OutputHeaders {
		planOutputHeaders[i] = elestio.LoadBalancerConfigOutputHeader{
			Key:   outputHeader.Key.ValueString(),
			Value: outputHeader.Value.ValueString(),
		}
	}

	planTargetServices := []string{}
	diags = plan.Config.TargetServices.ElementsAs(ctx, &planTargetServices, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	planRemoveResponseHeaders := []string{}
	diags = plan.Config.RemoveResponseHeaders.ElementsAs(ctx, &planRemoveResponseHeaders, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	payload := elestio.UpdateLoadBalancerConfigRequest{
		HostHeader:             plan.Config.HostHeader.ValueString(),
		IsAccessLogsEnabled:    plan.Config.AccessLogsEnabled.ValueBool(),
		IsForceHTTPSEnabled:    plan.Config.ForceHTTPSEnabled.ValueBool(),
		IPRateLimit:            plan.Config.IPRateLimitPerSecond.ValueInt64(),
		IsIPRateLimitEnabled:   plan.Config.IPRateLimitEnabled.ValueBool(),
		OutputCacheInSeconds:   plan.Config.OutputCacheInSeconds.ValueInt64(),
		IsStickySessionEnabled: plan.Config.StickySessionEnabled.ValueBool(),
		IsProxyProtocolEnabled: plan.Config.ProxyProtocolEnabled.ValueBool(),
		SSLDomains:             planSSLDomains,
		ForwardRules:           planForwardRules,
		OutputHeaders:          planOutputHeaders,
		TargetServices:         planTargetServices,
		RemoveResponseHeaders:  planRemoveResponseHeaders,
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
		Delay:                     30 * time.Second,
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
	loadBalancerModel.Config.AccessLogsEnabled = types.BoolValue(loadBalancerClient.Config.IsAccessLogsEnabled)
	loadBalancerModel.Config.ForceHTTPSEnabled = types.BoolValue(loadBalancerClient.Config.IsForceHTTPSEnabled)
	loadBalancerModel.Config.IPRateLimitPerSecond = types.Int64Value(loadBalancerClient.Config.IPRateLimit)
	loadBalancerModel.Config.IPRateLimitEnabled = types.BoolValue(loadBalancerClient.Config.IsIPRateLimitEnabled)
	loadBalancerModel.Config.OutputCacheInSeconds = types.Int64Value(loadBalancerClient.Config.OutputCacheInSeconds)
	loadBalancerModel.Config.StickySessionEnabled = types.BoolValue(loadBalancerClient.Config.IsStickySessionEnabled)
	loadBalancerModel.Config.ProxyProtocolEnabled = types.BoolValue(loadBalancerClient.Config.IsProxyProtocolEnabled)
	loadBalancerModel.Config.SSLDomains, diags = types.SetValueFrom(ctx, types.StringType, loadBalancerClient.Config.SSLDomains)
	if diags.HasError() {
		return loadBalancerModel, diags
	}
	forwardRules := make([]ForwardRuleModel, len(loadBalancerClient.Config.ForwardRules))
	for i, r := range loadBalancerClient.Config.ForwardRules {
		forwardRules[i] = ForwardRuleModel{
			Port:           types.StringValue(r.Port),
			Protocol:       types.StringValue(r.Protocol),
			TargetPort:     types.StringValue(r.TargetPort),
			TargetProtocol: types.StringValue(r.TargetProtocol),
		}
	}
	loadBalancerModel.Config.ForwardRules = forwardRules
	outputHeaders := make([]OutputHeaderModel, len(loadBalancerClient.Config.OutputHeaders))
	for i, h := range loadBalancerClient.Config.OutputHeaders {
		outputHeaders[i] = OutputHeaderModel{
			Key:   types.StringValue(h.Key),
			Value: types.StringValue(h.Value),
		}
	}
	loadBalancerModel.Config.OutputHeaders = outputHeaders
	loadBalancerModel.Config.TargetServices, diags = types.SetValueFrom(ctx, types.StringType, loadBalancerClient.Config.TargetServices)
	if diags.HasError() {
		return loadBalancerModel, diags
	}
	loadBalancerModel.Config.RemoveResponseHeaders, diags = types.SetValueFrom(ctx, types.StringType, loadBalancerClient.Config.RemoveResponseHeaders)
	if diags.HasError() {
		return loadBalancerModel, diags
	}
	loadBalancerModel.CreatedAt = types.StringValue(loadBalancerClient.CreatedAt)
	loadBalancerModel.CreatedBy = types.StringValue(loadBalancerClient.CreatorName)
	loadBalancerModel.DeploymentStatus = types.StringValue(loadBalancerClient.DeploymentStatus)
	loadBalancerModel.IPV4 = types.StringValue(loadBalancerClient.IPV4)
	loadBalancerModel.IPV6 = types.StringValue(loadBalancerClient.IPV6)
	loadBalancerModel.CNAME = types.StringValue(loadBalancerClient.CNAME)
	loadBalancerModel.Country = types.StringValue(loadBalancerClient.Country)
	loadBalancerModel.City = types.StringValue(loadBalancerClient.City)
	loadBalancerModel.GlobalIP = types.StringValue(loadBalancerClient.GlobalIP)
	loadBalancerModel.Cores = types.Int64Value(loadBalancerClient.Cores)
	loadBalancerModel.RAMSizeGB = types.StringValue(loadBalancerClient.RAMSizeGB)
	loadBalancerModel.StorageSizeGB = types.Int64Value(loadBalancerClient.StorageSizeGB)
	loadBalancerModel.PricePerHour = types.StringValue(loadBalancerClient.PricePerHour)
	return loadBalancerModel, diags
}
