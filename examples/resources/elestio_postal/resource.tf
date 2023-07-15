# Create and manage Postal service.
resource "elestio_postal" "demo_postal" {
  project_id    = "2500"
  server_name   = "demo-postal"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
