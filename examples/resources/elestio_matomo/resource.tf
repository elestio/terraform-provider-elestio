# Create and manage Matomo service.
resource "elestio_matomo" "demo_matomo" {
  project_id    = "2500"
  server_name   = "demo-matomo"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
