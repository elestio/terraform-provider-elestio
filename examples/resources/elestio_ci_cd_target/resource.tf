resource "elestio_ci_cd_target" "example" {
  project_id    = "2500"
  version       = "latest"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

}
