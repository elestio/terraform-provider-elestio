---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "elestio_centrifugo Resource - terraform-provider-elestio"
subcategory: ""
description: |-
  Centrifugo is a scalable real-time, language-agnostic messaging server. elestio_centrifugo is a preconfigured elestioservice resource (template_id: 151) running Centrifugo ([`dockerimage: centrifugo/centrifugo`](https://hub.docker.com/r/centrifugo/centrifugo)).
---

# elestio_centrifugo (Resource)

<img src="https://cf.appdrag.com/dashboard-openvm-clo-b2d42c/uploads/centrifugo-zHJH.svg" width="100" /><br/> Centrifugo is a scalable real-time, language-agnostic messaging server. <br/><br/>**elestio_centrifugo** is a preconfigured elestio_service resource (`template_id: 151`) running **Centrifugo** ([`docker_image: centrifugo/centrifugo`](https://hub.docker.com/r/centrifugo/centrifugo)).

## Example Usage

```terraform
# Create and manage Centrifugo service.
resource "elestio_centrifugo" "my_centrifugo" {
  project_id    = "2500"
  server_name   = "awesome-centrifugo"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `admin_email` (String) Service admin email. Requires replace to change it.
- `datacenter` (String) The datacenter of the provider where the service will be hosted. You can look for available datacenters in the [providers documentation](https://docs.elest.io/books/elestio-terraform-provider/page/providers-datacenters-and-server-types). Requires replace to change it.
- `project_id` (String) Identifier of the project in which the service is. Requires replace to change it.
- `provider_name` (String) The name of the provider to use to host the service. You can look for available provider names in the [providers documentation](https://docs.elest.io/books/elestio-terraform-provider/page/providers-datacenters-and-server-types). Requires replace to change it.
- `server_name` (String) Service server name. Must consist of lowercase letters, `a-z`, `0-9`, and `-`, and have a maximum length of 60 - underscore not allowed characters. Must be unique within the project. Requires replace to change it.
- `server_type` (String) The server type defines the power and memory allocated to the service. Each `provider_name` has a list of available server types. You can look for available server types in the [providers documentation](https://docs.elest.io/books/elestio-terraform-provider/page/providers-datacenters-and-server-types). You can only upgrade it, not downgrade.
- `ssh_keys` (Attributes Set) Indicate the list of SSH keys to add to the service. (see [below for nested schema](#nestedatt--ssh_keys))
- `support_level` (String) Service support level. You can look for available support levels and their advantages in the [pricing documentation](https://elest.io/pricing). Requires replace to change it in terraform. It is recommended to use the web dashboard to change it without replacing the service.

### Optional

- `alerts_enabled` (Boolean) Service alerts state. **Default** `true`.
- `app_auto_updates_enabled` (Boolean) Service app auto update state. **Default** `true`.
- `backups_enabled` (Boolean) Service backups state.  Requires a support_level higher than `level1`. **Default** `false`.
- `custom_domain_names` (Set of String) Indicate the list of domains for which you want to activate HTTPS / TLS / SSL. You will also need to create a DNS entry on your domain name (from your registrar control panel) pointing to your service. You must create a CNAME record pointing to the service `cname` value. Alternatively, you can create an A record pointing to the service `ipv4` value.
- `firewall_enabled` (Boolean) Service firewall state. **Default** `true`.
- `keep_backups_on_delete_enabled` (Boolean) Creates a backup and keeps all existing ones after deleting the service. If the project is deleted, the backups will be lost. **Default** `true`.
- `remote_backups_enabled` (Boolean) Service remote backups state. **Default** `true`.
- `system_auto_updates_enabled` (Boolean) Service system auto update state. **Default** `true`.
- `system_auto_updates_security_patches_only_enabled` (Boolean) Service system auto update security patches only state. **Default** `false`.
- `version` (String) This is the version of the software used as service. **Default** `latest`.

### Read-Only

- `admin` (Attributes) Service admin. (see [below for nested schema](#nestedatt--admin))
- `admin_user` (String) Service admin user.
- `app_auto_updates_day_of_week` (Number) Service app auto update day of week. `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`
- `app_auto_updates_hour` (Number) Service app auto update hour.
- `app_auto_updates_minute` (Number) Service app auto update minute.
- `category` (String) Service category.
- `city` (String) Service city.
- `cname` (String) Service CNAME.
- `cores` (Number) Service cores.
- `country` (String) Service country.
- `created_at` (String) Service creation date.
- `creator_name` (String) Service creator name.
- `database_admin` (Attributes) Service database admin. (see [below for nested schema](#nestedatt--database_admin))
- `deployment_ended_at` (String) Service deployment endedAt date.
- `deployment_started_at` (String) Service deployment startedAt date.
- `deployment_status` (String) Service deployement status.
- `env` (Map of String, Sensitive) Service environment variables.
- `external_backups_enabled` (Boolean) Service external backups state. **Default** `false`.
- `external_backups_retain_day_of_week` (Number) Service external backups retain day of week. `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`
- `external_backups_update_day_of_week` (Number) Service external backups update day. `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`
- `external_backups_update_hour` (Number) Service external backups update hour.
- `external_backups_update_minute` (Number) Service external backups update minute.
- `external_backups_update_type` (String) Service external backups update type.
- `firewall_id` (String) Service firewall id.
- `firewall_ports` (String) Service firewall ports.
- `global_ip` (String) Service global IP.
- `id` (String) Service identifier.
- `ipv4` (String) Service IPv4.
- `ipv6` (String) Service IPv6.
- `last_updated` (String)
- `price_per_hour` (String) Service price per hour.
- `ram_size_gb` (String) Service ram size in GB.
- `root_app_path` (String) Service root app path.
- `status` (String) Service status.
- `storage_size_gb` (Number) Service storage size in GB.
- `system_auto_updates_reboot_day_of_week` (Number) Service system auto update reboot day of week. `0 = Sunday`, `1 = Monday`, ..., `6 = Saturday`, `-1 = Everyday`
- `system_auto_updates_reboot_hour` (Number) Service system auto update reboot hour.
- `system_auto_updates_reboot_minute` (Number) Service system auto update reboot minute.
- `template_id` (Number) The template identifier defines the software used. You can look for available template ids in the [templates documentation](https://elest.io/fully-managed-services).
- `traffic_included` (Number) Service traffic included.
- `traffic_incoming` (Number) Service traffic incoming.
- `traffic_outgoing` (Number) Service traffic outgoing.

<a id="nestedatt--ssh_keys"></a>
### Nested Schema for `ssh_keys`

Required:

- `key_name` (String) SSH Key Name.
- `public_key` (String) SSH Public Key. With or without comment at the end. Example: `ssh-rsa AAAAB3Nz` or `ssh-rsa AAAAB3Nz comment@macbook.`


<a id="nestedatt--admin"></a>
### Nested Schema for `admin`

Read-Only:

- `password` (String, Sensitive) Service admin password.
- `url` (String) Service admin URL.
- `user` (String) Service admin user.


<a id="nestedatt--database_admin"></a>
### Nested Schema for `database_admin`

Read-Only:

- `command` (String, Sensitive) Service database admin command.
- `host` (String) Service database admin host.
- `password` (String, Sensitive) Service database admin password.
- `port` (String) Service database admin port.
- `user` (String) Service database admin user.

## Import

Import is supported using the following syntax:

```shell
# Import a Centrifugo service by specifying the Project ID it belongs to, and the Centrifugo service ID (spaced by a comma).
terraform import elestio_centrifugo.my_centrifugo project_id,centrifugo_service_id
```
