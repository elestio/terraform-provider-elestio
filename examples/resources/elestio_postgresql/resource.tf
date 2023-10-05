resource "elestio_postgresql" "example" {
  project_id    = "2500"
  version       = "16"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
