# Create and manage InfluxDB service.
resource "elestio_influxdb" "my_influxdb" {
  project_id    = "2500"
  server_name   = "awesome-influxdb"
  server_type   = "SMALL-1C-2G"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
