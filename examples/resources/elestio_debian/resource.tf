# Create and manage Debian service.
resource "elestio_debian" "my_debian" {
  project_id    = "2500"
  server_name   = "awesome-debian"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
