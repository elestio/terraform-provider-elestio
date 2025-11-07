resource "elestio_jitsi" "example" {
  project_id    = "2500"
  version       = "stable"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

}
