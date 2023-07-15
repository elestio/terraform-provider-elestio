# Create and manage Directus service.
resource "elestio_directus" "demo_directus" {
  project_id    = "2500"
  server_name   = "demo-directus"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
