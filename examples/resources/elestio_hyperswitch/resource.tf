resource "elestio_hyperswitch" "example" {
  project_id    = "2500"
  version       = "standalone"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}