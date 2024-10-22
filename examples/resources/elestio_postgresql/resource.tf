resource "elestio_postgresql" "example" {
  project_id    = "2500"
  version       = "16"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
