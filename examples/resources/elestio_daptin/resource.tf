# Create and manage Daptin service.
resource "elestio_daptin" "my_daptin" {
  project_id    = "2500"
  server_name   = "awesome-daptin"
  server_type   = "SMALL-1C-2G"
  version       = "master"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
