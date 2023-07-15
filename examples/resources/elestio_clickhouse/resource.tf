# Create and manage ClickHouse service.
resource "elestio_clickhouse" "demo_clickhouse" {
  project_id    = "2500"
  server_name   = "demo-clickhouse"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
