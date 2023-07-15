# Create and manage NextCloud service.
resource "elestio_nextcloud" "demo_nextcloud" {
  project_id    = "2500"
  server_name   = "demo-nextcloud"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
