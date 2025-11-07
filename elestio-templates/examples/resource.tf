resource "elestio_[TEMPLATE_RESOURCE_NAME]" "example" {
  project_id    = "2500"
  version       = "[TEMPLATE_DEFAULT_VERSION]"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

[TEMPLATE_FIREWALL_RULES]
}
