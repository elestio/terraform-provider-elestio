# Create and manage Wekan service.
resource "elestio_wekan" "demo_wekan" {
  project_id    = "2500"
  server_name   = "demo-wekan"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
