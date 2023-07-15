# Create and manage KeeWeb service.
resource "elestio_keeweb" "demo_keeweb" {
  project_id    = "2500"
  server_name   = "demo-keeweb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
