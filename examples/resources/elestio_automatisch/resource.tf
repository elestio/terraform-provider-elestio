# Create and manage Automatisch service.
resource "elestio_automatisch" "demo_automatisch" {
  project_id    = "2500"
  server_name   = "demo-automatisch"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
