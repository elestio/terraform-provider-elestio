# Create and manage Canopsis service.
resource "elestio_canopsis" "demo_canopsis" {
  project_id    = "2500"
  server_name   = "demo-canopsis"
  version       = "4.3.9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
