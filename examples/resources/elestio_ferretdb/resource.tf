resource "elestio_ferretdb" "example" {
  project_id    = "2500"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
