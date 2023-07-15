# Create and manage Shlink service.
resource "elestio_shlink" "demo_shlink" {
  project_id    = "2500"
  server_name   = "demo-shlink"
  version       = "stable"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
