# Create and manage Ubuntu service.
resource "elestio_ubuntu" "my_ubuntu" {
  project_id    = "2500"
  server_name   = "awesome-ubuntu"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
