# Create and manage TimescaleDB service.
resource "elestio_timescaledb" "my_timescaledb" {
  project_id    = "2500"
  server_name   = "awesome-timescaledb"
  server_type   = "SMALL-1C-2G"
  version       = "pg14-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
