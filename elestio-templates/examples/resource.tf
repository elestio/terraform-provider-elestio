resource "elestio_[TEMPLATE_RESOURCE_NAME]" "example" {
  project_id    = "2500"
  version       = "[TEMPLATE_DEFAULT_VERSION]"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
