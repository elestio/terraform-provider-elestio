# Create and manage Owncast service.
resource "elestio_owncast" "my_owncast" {
  project_id    = "2500"
  server_name   = "awesome-owncast"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
