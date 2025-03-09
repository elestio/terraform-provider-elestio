resource "elestio_tracardi" "example" {
  project_id    = "2500"
  version       = "1.0.2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
