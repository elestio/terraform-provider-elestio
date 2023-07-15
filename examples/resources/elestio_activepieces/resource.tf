# Create and manage Activepieces service.
resource "elestio_activepieces" "demo_activepieces" {
  project_id    = "2500"
  server_name   = "demo-activepieces"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
