resource "elestio_dittofeed" "example" {
  project_id    = "2500"
  version       = "v0.13.10"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
