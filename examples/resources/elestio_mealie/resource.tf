resource "elestio_mealie" "example" {
  project_id    = "2500"
  version       = "omni-nightly"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
