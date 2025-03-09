resource "elestio_hyperswitch" "example" {
  project_id    = "2500"
  version       = "standalone"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
