# Create and manage MSSQL service.
resource "elestio_mssql" "demo_mssql" {
  project_id    = "2500"
  server_name   = "demo-mssql"
  version       = "2019-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
