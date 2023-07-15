# Create and manage InfluxDB service.
resource "elestio_influxdb" "demo_influxdb" {
  project_id    = "2500"
  server_name   = "demo-influxdb"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
