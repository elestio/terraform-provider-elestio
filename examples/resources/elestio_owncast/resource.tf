# Create and manage Owncast service.
resource "elestio_owncast" "demo_owncast" {
  project_id    = "2500"
  server_name   = "demo-owncast"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
