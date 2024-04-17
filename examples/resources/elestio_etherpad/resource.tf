resource "elestio_etherpad" "example" {
  project_id    = "2500"
  version       = "v1.9.7"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
