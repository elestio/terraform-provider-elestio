resource "elestio_firefish" "example" {
  project_id    = "2500"
  version       = "v1.0.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
