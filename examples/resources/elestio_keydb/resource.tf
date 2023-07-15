# Create and manage KeyDB service.
resource "elestio_keydb" "demo_keydb" {
  project_id    = "2500"
  server_name   = "demo-keydb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
