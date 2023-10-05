resource "elestio_lago" "example" {
  project_id    = "2500"
  version       = "v0.48.0-beta"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
