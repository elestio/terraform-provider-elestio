resource "elestio_corteza" "example" {
  project_id    = "2500"
  version       = "2023.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
