# Create and manage Joplin service.
resource "elestio_joplin" "demo_joplin" {
  project_id    = "2500"
  server_name   = "demo-joplin"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
