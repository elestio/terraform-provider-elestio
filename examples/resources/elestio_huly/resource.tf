resource "elestio_huly" "example" {
  project_id    = "2500"
  version       = "v0.6.295"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
