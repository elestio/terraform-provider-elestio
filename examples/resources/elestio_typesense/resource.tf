resource "elestio_typesense" "example" {
  project_id    = "2500"
  version       = "27.0.rc34"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
