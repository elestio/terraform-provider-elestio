# Create and manage [TEMPLATE_DOCUMENTATION_NAME] service.
resource "elestio_[TEMPLATE_RESOURCE_NAME]" "my_[TEMPLATE_RESOURCE_NAME]" {
  project_id    = "2500"
  server_name   = "awesome-[TEMPLATE_RESOURCE_NAME]"
  server_type   = "SMALL-1C-2G"
  version       = "[TEMPLATE_DEFAULT_VERSION]"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
