# Create and manage Fider service.
resource "elestio_fider" "demo_fider" {
  project_id    = "2500"
  server_name   = "demo-fider"
  version       = "stable"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
