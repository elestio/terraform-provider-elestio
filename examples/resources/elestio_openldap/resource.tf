# Create and manage OpenLDAP service.
resource "elestio_openldap" "demo_openldap" {
  project_id    = "2500"
  server_name   = "demo-openldap"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
