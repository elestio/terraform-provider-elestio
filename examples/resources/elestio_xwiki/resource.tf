resource "elestio_xwiki" "example" {
  project_id    = "2500"
  version       = "16.4.0-postgres-tomcat"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
