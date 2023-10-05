resource "elestio_airbyte" "example" {
  project_id    = "2500"
  version       = "0.50.30"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
