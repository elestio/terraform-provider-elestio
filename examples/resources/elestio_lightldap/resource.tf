resource "elestio_lightldap" "example" {
  project_id    = "2500"
  version       = "stable"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
