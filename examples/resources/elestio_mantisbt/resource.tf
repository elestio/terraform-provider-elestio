# Create and manage MantisBT service.
resource "elestio_mantisbt" "demo_mantisbt" {
  project_id    = "2500"
  server_name   = "demo-mantisbt"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
