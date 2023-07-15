# Create and manage Mapzy service.
resource "elestio_mapzy" "demo_mapzy" {
  project_id    = "2500"
  server_name   = "demo-mapzy"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
