# Create and manage k0s service.
resource "elestio_k0s" "my_k0s" {
  project_id    = "2500"
  server_name   = "awesome-k0s"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
