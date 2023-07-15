# Create and manage OpenProject service.
resource "elestio_openproject" "demo_openproject" {
  project_id    = "2500"
  server_name   = "demo-openproject"
  version       = "12"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
