# Create and manage Superset service.
resource "elestio_superset" "demo_superset" {
  project_id    = "2500"
  server_name   = "demo-superset"
  version       = "master"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
