# Create and manage Authentik service.
resource "elestio_authentik" "my_authentik" {
  project_id    = "2500"
  server_name   = "awesome-authentik"
  server_type   = "SMALL-1C-2G"
  version       = "2021.12.5"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
