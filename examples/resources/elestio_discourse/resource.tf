resource "elestio_discourse" "example" {
  project_id    = "2500"
  version       = "3.4.7-debian-12-r0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
