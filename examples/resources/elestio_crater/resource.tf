resource "elestio_crater" "example" {
  project_id    = "2500"
  version       = "php8.1-deprecated"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
