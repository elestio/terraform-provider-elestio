# Create and manage TimescaleDB service.
resource "elestio_timescaledb" "demo_timescaledb" {
  project_id    = "2500"
  server_name   = "demo-timescaledb"
  version       = "pg14-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
