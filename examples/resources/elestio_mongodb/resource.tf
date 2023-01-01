# Create and manage MongoDB service.
resource "elestio_mongodb" "my_mongodb" {
  project_id    = "2500"
  server_name   = "awesome-mongodb"
  server_type   = "SMALL-1C-2G"
  version       = "6"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
