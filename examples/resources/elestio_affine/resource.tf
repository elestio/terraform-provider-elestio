# Create and manage Affine service.
resource "elestio_affine" "my_affine" {
  project_id    = "2500"
  server_name   = "awesome-affine"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
