# Create and manage OpenResty service.
resource "elestio_openresty" "my_openresty" {
  project_id    = "2500"
  server_name   = "awesome-openresty"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
