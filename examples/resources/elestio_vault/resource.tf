resource "elestio_vault" "example" {
  project_id    = "2500"
  version       = "1.13.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
