resource "elestio_mailu" "example" {
  project_id    = "2500"
  version       = "2.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
