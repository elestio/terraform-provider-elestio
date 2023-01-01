# Create and manage MySQL service.
resource "elestio_mysql" "my_mysql" {
  project_id    = "2500"
  server_name   = "awesome-mysql"
  server_type   = "SMALL-1C-2G"
  version       = "8"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
