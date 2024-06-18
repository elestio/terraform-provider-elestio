resource "elestio_huly" "example" {
  project_id    = "2500"
  version       = "v0.6.228a"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
