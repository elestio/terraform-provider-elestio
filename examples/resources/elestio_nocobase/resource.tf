# Create and manage NocoBase service.
resource "elestio_nocobase" "demo_nocobase" {
  project_id    = "2500"
  server_name   = "demo-nocobase"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
