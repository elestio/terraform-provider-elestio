# Create and manage Plausible Analytics service.
resource "elestio_plausible" "demo_plausible" {
  project_id    = "2500"
  server_name   = "demo-plausible"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
