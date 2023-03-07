# Create and manage Nebula service.
resource "elestio_nebula" "my_nebula" {
  project_id    = "2500"
  server_name   = "awesome-nebula"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
