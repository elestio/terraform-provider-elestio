# Create and manage Kimai service.
resource "elestio_kimai" "my_kimai" {
  project_id    = "2500"
  server_name   = "awesome-kimai"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
