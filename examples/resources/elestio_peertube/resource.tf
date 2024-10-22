resource "elestio_peertube" "example" {
  project_id    = "2500"
  version       = "production-bullseye"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
