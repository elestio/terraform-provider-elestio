# Create and manage Budibase service.
resource "elestio_budibase" "my_budibase" {
  project_id    = "2500"
  server_name   = "awesome-budibase"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
