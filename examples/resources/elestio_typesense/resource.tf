resource "elestio_typesense" "example" {
  project_id    = "2500"
  version       = "0.23.0.rc66"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
