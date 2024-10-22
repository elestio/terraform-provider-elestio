resource "elestio_dify" "example" {
  project_id    = "2500"
  version       = "main"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
