resource "elestio_lemmy" "example" {
  project_id    = "2500"
  version       = "0.19.4"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
