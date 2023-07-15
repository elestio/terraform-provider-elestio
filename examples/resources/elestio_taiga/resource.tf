# Create and manage Taiga service.
resource "elestio_taiga" "demo_taiga" {
  project_id    = "2500"
  server_name   = "demo-taiga"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
