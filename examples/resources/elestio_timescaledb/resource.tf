resource "elestio_timescaledb" "example" {
  project_id    = "2500"
  version       = "latest-pg17"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
