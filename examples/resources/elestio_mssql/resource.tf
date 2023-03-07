# Create and manage MSSQL service.
resource "elestio_mssql" "my_mssql" {
  project_id    = "2500"
  server_name   = "awesome-mssql"
  server_type   = "SMALL-1C-2G"
  version       = "2019-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
