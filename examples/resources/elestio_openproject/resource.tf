resource "elestio_openproject" "example" {
  project_id    = "2500"
  version       = "12"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
