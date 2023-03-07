# Create and manage Microk8s service.
resource "elestio_microk8s" "my_microk8s" {
  project_id    = "2500"
  server_name   = "awesome-microk8s"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
