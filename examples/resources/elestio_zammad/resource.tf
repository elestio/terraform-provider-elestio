resource "elestio_zammad" "example" {
  project_id    = "2500"
  version       = "5.4.1-38"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
