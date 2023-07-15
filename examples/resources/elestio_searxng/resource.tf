# Create and manage SearXNG service.
resource "elestio_searxng" "demo_searxng" {
  project_id    = "2500"
  server_name   = "demo-searxng"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
