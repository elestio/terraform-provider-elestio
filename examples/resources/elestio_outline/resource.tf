# Create and manage Outline service.
resource "elestio_outline" "demo_outline" {
  project_id    = "2500"
  server_name   = "demo-outline"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
