package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/elestio/elestio-go-api-client"
	"github.com/elestio/terraform-provider-elestio/internal/modifiers"
	"github.com/elestio/terraform-provider-elestio/internal/utils"
	"github.com/elestio/terraform-provider-elestio/internal/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

var (
	_ resource.Resource                   = &ServiceResource{}
	_ resource.ResourceWithValidateConfig = &ServiceResource{}
	_ resource.ResourceWithConfigure      = &ServiceResource{}
	_ resource.ResourceWithImportState    = &ServiceResource{}
)

type (
	ServiceTemplate struct {
		TemplateId         int64
		ResourceName       string
		DocumentationName  string
		Description        string
		DeprecationMessage string
		Logo               string
		DockerHubImage     string
		DefaultVersion     string
		Category           string
		FirewallPorts      []elestio.ServiceFirewallPort
	}

	ServiceResource struct {
		client *elestio.Client
		*ServiceTemplate
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
		DefaultPassword                             types.String `tfsdk:"default_password"`
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
		CustomDomainNames                           types.Set    `tfsdk:"custom_domain_names"`
		SSHKeys                                     types.Set    `tfsdk:"ssh_keys"`
		SSHPublicKeys                               types.Set    `tfsdk:"ssh_public_keys"`
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
		KeepBackupsOnDeleteEnabled                  types.Bool   `tfsdk:"keep_backups_on_delete_enabled"`
		FirewallEnabled                             types.Bool   `tfsdk:"firewall_enabled"`
		FirewallId                                  types.String `tfsdk:"firewall_id"`
		FirewallPorts                               types.String `tfsdk:"firewall_ports"`
		AlertsEnabled                               types.Bool   `tfsdk:"alerts_enabled"`
		LastUpdated                                 types.String `tfsdk:"last_updated"`
	}
)

type SSHKeyModel struct {
	// Deprecated: replaced by SSHPublicKeyModel
	KeyName   types.String `tfsdk:"key_name"`
	PublicKey types.String `tfsdk:"public_key"`
}

var sshKeyAttryTypes = map[string]attr.Type{
	// Deprecated: replaced by sshPublicKeyAttryTypes
	"key_name":   types.StringType,
	"public_key": types.StringType,
}

type AdminModel struct {
	URL      types.String `tfsdk:"url"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
}

var adminAttryTypes = map[string]attr.Type{
	"url":      types.StringType,
	"user":     types.StringType,
	"password": types.StringType,
}

type DatabaseAdminModel struct {
	Host     types.String `tfsdk:"host"`
	Port     types.String `tfsdk:"port"`
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
	Command  types.String `tfsdk:"command"`
}

var databaseAdminAttryTypes = map[string]attr.Type{
	"host":     types.StringType,
	"port":     types.StringType,
	"user":     types.StringType,
	"password": types.StringType,
	"command":  types.StringType,
}

func NewServiceResource(template *ServiceTemplate) resource.Resource {
	return &ServiceResource{
		ServiceTemplate: template,
	}
}

func (r *ServiceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.ResourceName
}

func (r *ServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	// ↓↓↓ Attributes that require a multi-stage construction process. ↓↓↓
	version := schema.StringAttribute{
		MarkdownDescription: "This is the version of the software used as service.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.RequiresReplaceIf(
				func(ctx context.Context, modifier planmodifier.StringRequest, resp *stringplanmodifier.RequiresReplaceIfFuncResponse) {
					if r.TemplateId == 11 {
						// PostgreSQL = templateId 11
						// PostgreSQL version cannot be upgraded
						resp.RequiresReplace = true
						return
					}

					// Add other templateId here if they cannot be upgraded

					resp.RequiresReplace = false
				},
				"This resource requires replace if you want to upgrade version.",
				"This resource Requires replace if you want to upgrade version.",
			),
		},
	}

	if r.DefaultVersion == "" {
		version.Required = true
	} else {
		// The version is optional only if the template specifies a default version
		version.MarkdownDescription += fmt.Sprintf(" **Default** `%s`.", r.DefaultVersion)
		version.Default = stringdefault.StaticString(r.DefaultVersion)
		version.Optional = true
		version.Computed = true
	}
	// ↑↑↑ Attributes that require a multi-stage construction process. ↑↑↑

	schemaMardownDescription := ""
	if r.TemplateId == 0 {
		schemaMardownDescription += "This resource is the generic way to create a service." +
			" You can choose the software by providing the `template_id` as a parameter." +
			" You can look for available template ids in the [templates documentation](https://elest.io/fully-managed-services)."
	} else {
		if r.Logo != "" {
			schemaMardownDescription += fmt.Sprintf(`<img src="%s" width="100" /><br/>`, r.Logo)
		}

		if r.Description != "" {
			schemaMardownDescription += fmt.Sprintf(" %s<br/><br/>", r.Description)
		}

		schemaMardownDescription += fmt.Sprintf("**elestio_%s** is a preconfigured elestio_service resource (`template_id: %d`) running **%s**", r.ResourceName, r.TemplateId, r.DocumentationName)

		if r.DockerHubImage != "" {
			schemaMardownDescription += fmt.Sprintf(" from the [Docker image](https://hub.docker.com/r/%s) `%s`", r.DockerHubImage, r.DockerHubImage)
		}

		schemaMardownDescription += "."
	}

	defaultSSHKeys, diags := types.SetValue(types.ObjectType{AttrTypes: sshKeyAttryTypes}, []attr.Value{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: schemaMardownDescription,
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
				MarkdownDescription: "Identifier of the project in which the service is." +
					" Requires replace to change it.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_name": schema.StringAttribute{
				MarkdownDescription: "Service server name." +
					" Must consist of lowercase letters, `a-z`, `0-9`, and `-`, and have a maximum length of 60 - underscore not allowed characters." +
					" Must be unique within the project." +
					" Requires replace to change it.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 60),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-z0-9-]+$`),
						"Must consist of lowercase letters, a-z, 0-9, and - (dash), _ (underscore) not allowed characters.",
					),
				},
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "The server type defines the power and memory allocated to the service." +
					" Each `provider_name` has a list of available server types." +
					" You can look for available server types in the [providers documentation](https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/providers_datacenters_server_types)." +
					" You can only upgrade it, not downgrade." +
					"<br/>Requires replace to update the server type with the provider `scale_way`.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(modifiers.RequiresReplaceIfProviderScaleway(),
						"Requires replace to update the server type with the provider `scale_way`.",
						"Requires replace to update the server type with the provider `scale_way`.",
					),
				},
			},
			"template_id": schema.Int64Attribute{
				MarkdownDescription: " The template identifier defines the software used." +
					" You can look for available template ids in the [templates documentation](https://elest.io/fully-managed-services).",
				Required: r.TemplateId == 0,
				Computed: r.TemplateId != 0,
				PlanModifiers: []planmodifier.Int64{
					utils.If(r.TemplateId == 0, int64planmodifier.RequiresReplace(), int64planmodifier.UseStateForUnknown()),
				},
			},
			"version": version,
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "The name of the provider to use to host the service." +
					" You can look for available provider names in the [providers documentation](https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/providers_datacenters_server_types)." +
					" Requires replace to change it.",
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
				MarkdownDescription: "The datacenter of the provider where the service will be hosted." +
					" You can look for available datacenters in the [providers documentation](https://registry.terraform.io/providers/elestio/elestio/latest/docs/guides/providers_datacenters_server_types)." +
					" Requires replace to change it.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"support_level": schema.StringAttribute{
				MarkdownDescription: "Service support level." +
					" Available support levels are `level1`, `level2` and `level3`." +
					" You can look for their advantages in the [pricing documentation](https://elest.io/pricing)." +
					" Requires replace the whole resource to change it in terraform." +
					" It is recommended to use the web dashboard to change it without replacing the service.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("level1"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("level1", "level2", "level3"),
				},
			},
			"admin_email": schema.StringAttribute{
				MarkdownDescription: "Service admin email." +
					" Requires replace to change it.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.IsEmail(),
				},
			},
			"default_password": schema.StringAttribute{
				MarkdownDescription: "Set the default password used by you services at **CREATION** time." +
					"</br>The password can only contain alphanumeric characters or hyphens `-`." +
					" Require at least 10 characters, one uppercase letter, one lowercase letter and one number." +
					"</br>If you don't set a password, a random one will be generated by the API." +
					"</br>This attribute will **not be synced** after the creation." +
					" Use `admin.password` or `database_admin.password` to get the current password after the creation.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					validators.IsDefaultPassword(),
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ipv6": schema.StringAttribute{
				MarkdownDescription: "Service IPv6.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cname": schema.StringAttribute{
				MarkdownDescription: "Service CNAME.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_domain_names": schema.SetAttribute{
				MarkdownDescription: "Indicate the list of domains for which you want to activate HTTPS / TLS / SSL." +
					" You will also need to create a DNS entry on your domain name (from your registrar control panel) pointing to your service." +
					" You must create a CNAME record pointing to the service `cname` value." +
					" Alternatively, you can create an A record pointing to the service `ipv4` value.",
				// Databases & Cache do not support custom domains
				Optional:    r.Category != "Databases & Cache",
				Computed:    true,
				ElementType: types.StringType,
			},
			"ssh_keys": schema.SetNestedAttribute{
				MarkdownDescription: "This attribute allows you to add SSH keys to your service.",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(defaultSSHKeys),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key_name": schema.StringAttribute{
							MarkdownDescription: "SSH Key Name.",
							Required:            true,
						},
						"public_key": schema.StringAttribute{
							MarkdownDescription: "SSH Public Key." +
								"The SSH public key should only contain two parts separated by a space." +
								" Example: `ssh-rsa AAaCfa...WAqDUNs=`." +
								" You should not include the username, hostname, or comment.",
							Required: true,
							Validators: []validator.String{
								validators.IsSSHPublicKey(),
							},
							DeprecationMessage: "This attribute is deprecated and will be removed in a future version." +
								" Please use `ssh_public_keys` instead.",
						},
					},
				},
			},
			"ssh_public_keys": sshPublicKeysSchema,
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
				Sensitive:           true,
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
					},
				},
			},
			"database_admin": schema.SingleNestedAttribute{
				MarkdownDescription: "Service database admin.",
				Computed:            true,
				Sensitive:           true,
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
					},
					"command": schema.StringAttribute{
						MarkdownDescription: "Service database admin command.",
						Computed:            true,
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
				MarkdownDescription: "Service ram size in GB.",
				Computed:            true,
			},
			"storage_size_gb": schema.Int64Attribute{
				MarkdownDescription: "Service storage size in GB.",
				Computed:            true,
			},
			"price_per_hour": schema.StringAttribute{
				MarkdownDescription: "Service price per hour.",
				Computed:            true,
			},
			"app_auto_updates_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service app auto update state." +
					" **Default** `true`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"app_auto_updates_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service app auto update day of week." +
					" `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`",
				Computed: true,
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
				MarkdownDescription: "Service system auto update state." +
					" **Default** `true`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"system_auto_updates_security_patches_only_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service system auto update security patches only state." +
					" **Default** `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"system_auto_updates_reboot_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service system auto update reboot day of week." +
					" `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`",
				Computed: true,
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
				MarkdownDescription: "Service backups state. " +
					" Requires a support_level higher than `level1`." +
					" **Default** `false`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"remote_backups_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service remote backups state." +
					" **Default** `true`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"keep_backups_on_delete_enabled": schema.BoolAttribute{
				MarkdownDescription: "Creates a backup and keeps all existing ones after deleting the service." +
					" If the project is deleted, the backups will be lost." +
					" **Default** `true`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"external_backups_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service external backups state." +
					" **Default** `false`.",
				Computed: true,
				// TODO: Handle external backups with s3 config
			},
			"external_backups_update_day_of_week": schema.Int64Attribute{
				MarkdownDescription: "Service external backups update day." +
					" `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`",
				Computed: true,
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
				MarkdownDescription: "Service external backups retain day of week." +
					" `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`",
				Computed: true,
			},
			"firewall_enabled": schema.BoolAttribute{
				MarkdownDescription: "Service firewall state." +
					" **Default** `true`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
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
				MarkdownDescription: "Service alerts state." +
					" **Default** `true`.",
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
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

	// TODO: Move this check to a validator file with the proper format
	ensureSSHPublicKeysUsernamesAreUnique(&ctx, &data.SSHPublicKeys, &resp.Diagnostics)
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

	service, err := r.createServiceWithRetry(ctx, elestio.CreateServiceRequest{
		ProjectID:    data.ProjectID.ValueString(),
		ServerName:   data.ServerName.ValueString(),
		ServerType:   data.ServerType.ValueString(),
		AppPassword:  data.DefaultPassword.ValueString(),
		TemplateID:   templateId,
		Version:      data.Version.ValueString(),
		ProviderName: data.ProviderName.ValueString(),
		Datacenter:   data.Datacenter.ValueString(),
		SupportLevel: data.SupportLevel.ValueString(),
		AdminEmail:   data.AdminEmail.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Service",
			fmt.Sprintf("Unable to create service, got error: %s", err),
		)
		return
	}

	// Service is created but we need to wait for the default configuration to be applied.
	serviceConfigured, err := r.waitServiceDefaultConfiguration(ctx, service)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Service",
			fmt.Sprintf("Unable to wait service default configuration, got error: %s", err),
		)
		return
	}

	data.Id = types.StringValue(serviceConfigured.ID)

	// Update some fields that are not available in the create request.
	serviceUpdated, err := r.updateElestioService(ctx, serviceConfigured, data, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Service",
			fmt.Sprintf("Unable after create to update fields that are not available in the create request, got error: %s", err),
		)
		return
	}

	convertElestioToTerraformFormat(ctx, data, serviceUpdated, &resp.Diagnostics)
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

	convertElestioToTerraformFormat(ctx, data, service, &resp.Diagnostics)
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

	convertElestioToTerraformFormat(ctx, plan, updatedService, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	projectId, serviceId, keepBackups := state.ProjectID.ValueString(), state.Id.ValueString(), state.KeepBackupsOnDeleteEnabled.ValueBool()
	service, err := r.client.Service.Get(projectId, serviceId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Service",
			fmt.Sprintf("Unable to get service, got error: %s", err),
		)
		return
	}

	if err := r.deleteServiceWithRetry(ctx, service.ProjectID, service.ID, keepBackups); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Service",
			fmt.Sprintf("Unable to start service deletion, got error: %s", err),
		)
		return
	}

	if err := r.waitServiceDeletion(ctx, service); err != nil {
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
	convertElestioToTerraformFormat(ctx, state, service, diags)

	// Server type update should be done first, because it requires to stop the service
	if !state.ServerType.Equal(plan.ServerType) {
		if err := r.client.Service.UpdateServerType(service.ID, plan.ServerType.ValueString(), service.ProviderName, service.Datacenter); err != nil {
			return nil, fmt.Errorf("failed to update serverType: %s", err)
		}
		r.waitServerReboot(ctx, service)
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
			if err := r.client.Service.EnableFirewall(service.ID, r.FirewallPorts); err != nil {
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

	if r.Category != "Databases & Cache" {
		var stateCustomDomainNames []string
		state.CustomDomainNames.ElementsAs(ctx, &stateCustomDomainNames, false)

		var planCustomDomainNames []string
		plan.CustomDomainNames.ElementsAs(ctx, &planCustomDomainNames, false)

		for _, planCustomDomainName := range planCustomDomainNames {
			if !utils.Contains(stateCustomDomainNames, planCustomDomainName) {
				if err := r.client.Service.AddCustomDomainName(service.ID, planCustomDomainName); err != nil {
					return nil, fmt.Errorf("failed to add customDomainName: %s", err)
				}
			}
		}

		for _, stateCustomDomainName := range stateCustomDomainNames {
			if !utils.Contains(planCustomDomainNames, stateCustomDomainName) {
				if err := r.client.Service.RemoveCustomDomainName(service.ID, stateCustomDomainName); err != nil {
					return nil, fmt.Errorf("failed to remove customDomainName: %s", err)
				}
			}
		}
	}

	keysToAdd, keysToUpdate, keysToRemove := compareSSHPublicKeys(&ctx, &state.SSHPublicKeys, &plan.SSHPublicKeys, diags)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to compare ssh public keys from state to plan")
	}

	for _, key := range keysToRemove {
		if err := r.client.Service.RemoveSSHPublicKey(service.ID, key.Username.ValueString()); err != nil {
			return nil, fmt.Errorf("failed to remove ssh public key: %s", err)
		}
	}

	for _, key := range keysToUpdate {
		if err := r.client.Service.RemoveSSHPublicKey(service.ID, key.Username.ValueString()); err != nil {
			return nil, fmt.Errorf("failed to update (remove the old one) ssh public key: %s", err)
		}
		if err := r.client.Service.AddSSHPublicKey(service.ID, key.Username.ValueString(), key.KeyData.ValueString()); err != nil {
			return nil, fmt.Errorf("failed to update (add the new one) ssh public key: %s", err)
		}
	}

	for _, key := range keysToAdd {
		if err := r.client.Service.AddSSHPublicKey(service.ID, key.Username.ValueString(), key.KeyData.ValueString()); err != nil {
			return nil, fmt.Errorf("failed to add ssh public key: %s", err)
		}
	}

	// Scaleway does not support updating ssh keys without rebooting the server.
	keyWasUpdated := len(keysToAdd) > 0 || len(keysToUpdate) > 0 || len(keysToRemove) > 0
	if keyWasUpdated && plan.ProviderName.ValueString() == "scaleway" {
		if err := r.client.Service.RebootServer(service.ID); err != nil {
			return nil, fmt.Errorf("failed to reboot server to update scaleway ssh keys: %s", err)
		}
		r.waitServerReboot(ctx, service)
	}

	service, err := r.client.Service.Get(service.ProjectID, service.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %s", err)
	}
	return service, nil
}

func convertElestioToTerraformFormat(ctx context.Context, data *ServiceResourceModel, service *elestio.Service, diags *diag.Diagnostics) {
	var d diag.Diagnostics

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
	data.CustomDomainNames = utils.SliceStringToSetType(service.CustomDomainNames, diags)
	sshPublicKeys := make([]SSHPublicKeyModel, len(service.SSHPublicKeys))
	for i, s := range service.SSHPublicKeys {
		sshPublicKeys[i] = SSHPublicKeyModel{
			Username: types.StringValue(s.Name),
			KeyData:  types.StringValue(s.Key),
		}
	}
	setSSHPublicKeys, d := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: sshPublicKeyAttryTypes}, sshPublicKeys)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	data.SSHPublicKeys = setSSHPublicKeys
	data.Country = types.StringValue(service.Country)
	data.City = types.StringValue(service.City)
	data.AdminUser = types.StringValue(service.AdminUser)
	data.RootAppPath = types.StringValue(service.RootAppPath)
	data.Env = utils.MapStringToMapType(service.Env, diags)
	data.Admin, d = types.ObjectValue(adminAttryTypes, map[string]attr.Value{
		"url":      types.StringValue(service.Admin.URL),
		"user":     types.StringValue(service.Admin.User),
		"password": types.StringValue(service.Admin.Password),
	})
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	data.DatabaseAdmin, d = types.ObjectValue(databaseAdminAttryTypes, map[string]attr.Value{
		"host":     types.StringValue(service.DatabaseAdmin.Host),
		"port":     types.StringValue(service.DatabaseAdmin.Port),
		"user":     types.StringValue(service.DatabaseAdmin.User),
		"password": types.StringValue(service.DatabaseAdmin.Password),
		"command":  types.StringValue(service.DatabaseAdmin.Command),
	})
	diags.Append(d...)
	if diags.HasError() {
		return
	}
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
	data.SystemAutoUpdatesSecurityPatchesOnlyEnabled = utils.If(
		// condition
		!data.SystemAutoUpdatesEnabled.ValueBool(),
		// if true
		types.BoolValue(false),
		// if false
		utils.BoolValue(service.SystemAutoUpdatesSecurityPatchesOnlyEnabled),
	)
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

func (r *ServiceResource) createServiceWithRetry(ctx context.Context, request elestio.CreateServiceRequest) (*elestio.Service, error) {
	var serviceR *elestio.Service
	var err error
	retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
		serviceR, err = r.client.Service.Create(request)
		if err != nil {
			return retry.RetryableError(err)
		}

		return nil
	})

	return serviceR, err
}

func (r *ServiceResource) deleteServiceWithRetry(ctx context.Context, projectId string, serviceId string, keepBackups bool) error {
	timeout := 2 * time.Minute
	stateConf := retry.StateChangeConf{
		Pending: []string{"error"},
		Target:  []string{"success"},
		Refresh: func() (interface{}, string, error) {
			err := r.client.Service.Delete(projectId, serviceId, keepBackups)
			if err != nil {
				// Retry on error
				// Do not return error here, because it will stop the loop -> return nil, "error", err
				return struct{}{}, "error", nil
			}

			return struct{}{}, "success", nil
		},
		Timeout:                   timeout,
		Delay:                     0,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 0,
	}

	tflog.Trace(ctx, fmt.Sprintf("DeleteServiceWithRetry timeout %.0f minutes", timeout.Minutes()))

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("DeleteServiceWithRetry failed, got error: %s", err)
	}

	return nil
}

func (r *ServiceResource) waitServiceDefaultConfiguration(ctx context.Context, service *elestio.Service) (*elestio.Service, error) {
	timeout := 15 * time.Minute
	stateConf := retry.StateChangeConf{
		Pending: []string{"waiting"},
		Target:  []string{"configured"},
		Refresh: func() (interface{}, string, error) {
			serviceR, err := r.client.Service.Get(service.ProjectID, service.ID)
			if err != nil {
				return struct{}{}, "", err
			}

			if serviceR.DeploymentStatus != elestio.ServiceDeploymentStatusDeployed {
				return struct{}{}, "waiting", nil
			}

			// App auto updates are enabled by default at service creation
			if !utils.BoolValue(serviceR.AppAutoUpdatesEnabled).ValueBool() {
				return struct{}{}, "waiting", nil
			}

			// System auto updates are enabled by default at service creation
			if !utils.BoolValue(serviceR.SystemAutoUpdatesEnabled).ValueBool() {
				return struct{}{}, "waiting", nil
			}

			// Backups are enabled by default at service creation if service level is greater than level1
			if serviceR.SupportLevel != "level1" && !utils.BoolValue(serviceR.BackupsEnabled).ValueBool() {
				return struct{}{}, "waiting", nil
			}

			// Remote backups are enabled by default at service creation
			if !utils.BoolValue(serviceR.RemoteBackupsEnabled).ValueBool() {
				return struct{}{}, "waiting", nil
			}

			// Firewall is enabled by default at service creation
			if !utils.BoolValue(serviceR.FirewallEnabled).ValueBool() {
				return struct{}{}, "waiting", nil
			}

			// Alerts are enabled by default at service creation
			if !utils.BoolValue(serviceR.AlertsEnabled).ValueBool() {
				return struct{}{}, "waiting", nil
			}

			return serviceR, "configured", nil
		},
		Timeout:                   timeout,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 3,
	}

	tflog.Trace(ctx, fmt.Sprintf("WaitServiceDefaultConfiguration timeout %.0f minutes", timeout.Minutes()))
	serviceConfigured, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("WaitServiceDefaultConfiguration failed, got error: %s", err)
	}

	return serviceConfigured.(*elestio.Service), nil
}

func (r *ServiceResource) waitServiceDeletion(ctx context.Context, service *elestio.Service) error {
	timeout := 15 * time.Minute
	stateConf := retry.StateChangeConf{
		Pending: []string{"waiting"},
		Target:  []string{"deleted"},
		Refresh: func() (any, string, error) {
			_, err := r.client.Service.Get(service.ProjectID, service.ID)

			if err == nil {
				return struct{}{}, "waiting", nil
			}

			return struct{}{}, "deleted", nil
		},
		Timeout:                   timeout,
		Delay:                     60 * time.Second,
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	tflog.Trace(ctx, fmt.Sprintf("WaitServiceDeletion timeout %.0f minutes", timeout.Minutes()))
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("WaitServiceDeletion failed, got error: %s", err)
	}

	return nil
}

func (r *ServiceResource) waitServerReboot(ctx context.Context, service *elestio.Service) (*elestio.Service, error) {
	timeout := 10 * time.Minute
	stateConf := retry.StateChangeConf{
		Pending: []string{"rebooting"},
		Target:  []string{"rebooted"},
		Refresh: func() (interface{}, string, error) {
			serviceW, err := r.client.Service.Get(service.ProjectID, service.ID)
			if err != nil {
				return struct{}{}, "", err
			}

			// running -> stopping -> stopped -> migrating -> running
			if serviceW.Status != elestio.ServiceStatusRunning {
				return struct{}{}, "rebooting", nil
			}

			return serviceW, "rebooted", nil
		},
		Timeout:                   timeout,
		Delay:                     40 * time.Second,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 2,
	}

	tflog.Trace(ctx, fmt.Sprintf("Service reboot waiter timeout %.0f minutes", timeout.Minutes()))

	serviceRebooted, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("service reboot waiter failed, got error: %s", err)
	}

	return serviceRebooted.(*elestio.Service), nil
}
