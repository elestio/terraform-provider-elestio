resource "elestio_twenty" "example" {
  project_id    = "2500"
  version       = "v0.30.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
