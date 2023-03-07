# Create and manage PostgreSQL service.
resource "elestio_postgresql" "my_postgresql" {
  project_id    = "2500"
  server_name   = "awesome-postgresql"
  server_type   = "SMALL-1C-2G"
  version       = "14"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
