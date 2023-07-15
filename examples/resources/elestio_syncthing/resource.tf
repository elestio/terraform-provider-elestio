# Create and manage Syncthing service.
resource "elestio_syncthing" "demo_syncthing" {
  project_id    = "2500"
  server_name   = "demo-syncthing"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
