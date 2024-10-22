resource "elestio_zammad" "example" {
  project_id    = "2500"
  version       = "6.2.0-1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
