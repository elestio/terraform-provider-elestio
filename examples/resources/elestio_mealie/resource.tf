# Create and manage Mealie service.
resource "elestio_mealie" "my_mealie" {
  project_id    = "2500"
  server_name   = "awesome-mealie"
  server_type   = "SMALL-1C-2G"
  version       = "omni-nightly"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
