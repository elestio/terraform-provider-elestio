resource "elestio_wekan" "example" {
  project_id    = "2500"
  version       = "v7.81"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
