# Create and manage K3S service.
resource "elestio_k3s" "my_k3s" {
  project_id    = "2500"
  server_name   = "awesome-k3s"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
