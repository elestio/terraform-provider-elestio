resource "elestio_filestash" "example" {
  project_id    = "2500"
  version       = "6b271d3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
