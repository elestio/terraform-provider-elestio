resource "elestio_chatwoot" "example" {
  project_id    = "2500"
  version       = "v3.6.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
