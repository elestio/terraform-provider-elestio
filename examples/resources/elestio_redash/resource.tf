resource "elestio_redash" "example" {
  project_id    = "2500"
  version       = "25.1.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
