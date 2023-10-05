resource "elestio_authentik" "example" {
  project_id    = "2500"
  version       = "2023.3.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
