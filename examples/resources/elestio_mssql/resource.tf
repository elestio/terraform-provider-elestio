resource "elestio_mssql" "example" {
  project_id    = "2500"
  version       = "2019-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
