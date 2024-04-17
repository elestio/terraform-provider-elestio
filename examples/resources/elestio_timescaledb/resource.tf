resource "elestio_timescaledb" "example" {
  project_id    = "2500"
  version       = "latest-pg16"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
