# Create and manage Neko Rooms service.
resource "elestio_neko" "my_neko" {
  project_id    = "2500"
  server_name   = "awesome-neko"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
