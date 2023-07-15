# Create and manage Budibase service.
resource "elestio_budibase" "demo_budibase" {
  project_id    = "2500"
  server_name   = "demo-budibase"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
