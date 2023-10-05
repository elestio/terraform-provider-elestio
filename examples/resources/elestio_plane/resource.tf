# Create and manage Plane service.
resource "elestio_plane" "demo_plane" {
  project_id    = "2500"
  server_name   = "demo-plane"
  version       = "0.11"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
