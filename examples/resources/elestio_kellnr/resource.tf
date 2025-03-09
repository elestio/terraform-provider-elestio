resource "elestio_kellnr" "example" {
  project_id    = "2500"
  version       = "5.0.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
