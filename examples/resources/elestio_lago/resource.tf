# Create and manage Lago service.
resource "elestio_lago" "demo_lago" {
  project_id    = "2500"
  server_name   = "demo-lago"
  version       = "v0.48.0-beta"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
