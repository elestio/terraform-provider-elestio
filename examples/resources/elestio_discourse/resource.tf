# Create and manage Discourse service.
resource "elestio_discourse" "demo_discourse" {
  project_id    = "2500"
  server_name   = "demo-discourse"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
