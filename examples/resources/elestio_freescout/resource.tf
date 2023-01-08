# Create and manage FreeScout service.
resource "elestio_freescout" "my_freescout" {
  project_id    = "2500"
  server_name   = "awesome-freescout"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
