# Create and manage MariaDB service.
resource "elestio_mariadb" "demo_mariadb" {
  project_id    = "2500"
  server_name   = "demo-mariadb"
  version       = "10.9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
