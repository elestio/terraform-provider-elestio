resource "elestio_tracardi" "example" {
  project_id    = "2500"
  version       = "0.8.0.5"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
