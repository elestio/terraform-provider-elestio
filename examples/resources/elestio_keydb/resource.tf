# Create and manage KeyDB service.
resource "elestio_keydb" "my_keydb" {
  project_id    = "2500"
  server_name   = "awesome-keydb"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
