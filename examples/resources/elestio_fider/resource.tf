# Create and manage Fider service.
resource "elestio_fider" "my_fider" {
  project_id    = "2500"
  server_name   = "awesome-fider"
  server_type   = "SMALL-1C-2G"
  version       = "stable"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
