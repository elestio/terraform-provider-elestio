# Create and manage Baserow service.
resource "elestio_baserow" "demo_baserow" {
  project_id    = "2500"
  server_name   = "demo-baserow"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
