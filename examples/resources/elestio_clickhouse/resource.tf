# Create and manage ClickHouse service.
resource "elestio_clickhouse" "my_clickhouse" {
  project_id    = "2500"
  server_name   = "awesome-clickhouse"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
