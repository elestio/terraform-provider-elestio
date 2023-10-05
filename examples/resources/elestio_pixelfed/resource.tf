# Create and manage Pixelfed service.
resource "elestio_pixelfed" "demo_pixelfed" {
  project_id    = "2500"
  server_name   = "demo-pixelfed"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
