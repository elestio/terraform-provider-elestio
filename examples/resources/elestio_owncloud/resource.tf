# Create and manage ownCloud service.
resource "elestio_owncloud" "demo_owncloud" {
  project_id    = "2500"
  server_name   = "demo-owncloud"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
