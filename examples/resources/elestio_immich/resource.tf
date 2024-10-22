resource "elestio_immich" "example" {
  project_id    = "2500"
  version       = "release"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
