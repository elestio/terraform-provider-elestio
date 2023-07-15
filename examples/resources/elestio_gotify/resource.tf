# Create and manage Gotify service.
resource "elestio_gotify" "demo_gotify" {
  project_id    = "2500"
  server_name   = "demo-gotify"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
