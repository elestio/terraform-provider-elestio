resource "elestio_friendica" "example" {
  project_id    = "2500"
  version       = "2023.05"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
