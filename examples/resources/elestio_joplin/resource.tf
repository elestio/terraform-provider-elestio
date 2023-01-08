# Create and manage Joplin service.
resource "elestio_joplin" "my_joplin" {
  project_id    = "2500"
  server_name   = "awesome-joplin"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
