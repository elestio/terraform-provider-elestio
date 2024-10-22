resource "elestio_agencyos" "example" {
  project_id    = "2500"
  version       = "v9.26.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
