resource "elestio_tracardi" "example" {
  project_id    = "2500"
  version       = "0.7.1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
