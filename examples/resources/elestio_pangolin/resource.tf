resource "elestio_pangolin" "example" {
  project_id    = "2500"
  version       = "1.4.0"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

}
