# Create and manage NextCloud service.
resource "elestio_nextcloud" "my_nextcloud" {
  project_id    = "2500"
  server_name   = "awesome-nextcloud"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
