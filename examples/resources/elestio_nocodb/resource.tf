# Create and manage NocoDB service.
resource "elestio_nocodb" "my_nocodb" {
  project_id    = "2500"
  server_name   = "awesome-nocodb"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
