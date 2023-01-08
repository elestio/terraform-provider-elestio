# Create and manage OpenProject service.
resource "elestio_openproject" "my_openproject" {
  project_id    = "2500"
  server_name   = "awesome-openproject"
  server_type   = "SMALL-1C-2G"
  version       = "12"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
