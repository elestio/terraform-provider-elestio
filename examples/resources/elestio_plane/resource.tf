resource "elestio_plane" "example" {
  project_id    = "2500"
  version       = "0.11"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
