resource "elestio_firefish" "example" {
  project_id    = "2500"
  version       = "v1.0.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
