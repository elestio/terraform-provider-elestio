resource "elestio_emqx" "example" {
  project_id    = "2500"
  version       = "5.5.1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
