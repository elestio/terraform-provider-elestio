# Create and manage [TEMPLATE_DOCUMENTATION_NAME] service.
resource "elestio_[TEMPLATE_RESOURCE_NAME]" "demo_[TEMPLATE_RESOURCE_NAME]" {
  project_id    = "2500"
  server_name   = "demo-[TEMPLATE_RESOURCE_NAME]"
  version       = "[TEMPLATE_DEFAULT_VERSION]"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
