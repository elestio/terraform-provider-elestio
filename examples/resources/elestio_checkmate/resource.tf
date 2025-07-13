resource "elestio_checkmate" "example" {
  project_id    = "2500"
  version       = "frontend-dist"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
