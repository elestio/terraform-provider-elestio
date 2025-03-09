resource "elestio_corteza" "example" {
  project_id    = "2500"
  version       = "2024.9.0-rc.5"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
