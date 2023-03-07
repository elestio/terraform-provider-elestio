# Create and manage Portainer service.
resource "elestio_portainer" "my_portainer" {
  project_id    = "2500"
  server_name   = "awesome-portainer"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
