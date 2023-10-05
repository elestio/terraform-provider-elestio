resource "elestio_rocket_chat" "example" {
  project_id    = "2500"
  version       = "sha-856c235"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
