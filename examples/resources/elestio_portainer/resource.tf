# Create and manage Portainer service.
resource "elestio_portainer" "demo_portainer" {
  project_id    = "2500"
  server_name   = "demo-portainer"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
