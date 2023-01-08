# Create and manage Squid service.
resource "elestio_squid" "my_squid" {
  project_id    = "2500"
  server_name   = "awesome-squid"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
