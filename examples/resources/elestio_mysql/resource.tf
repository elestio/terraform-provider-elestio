# Create and manage MySQL service.
resource "elestio_mysql" "demo_mysql" {
  project_id    = "2500"
  server_name   = "demo-mysql"
  version       = "8"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
