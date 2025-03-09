resource "elestio_tooljet" "example" {
  project_id    = "2500"
  version       = "ce-lts-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
