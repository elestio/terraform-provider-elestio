# Create and manage LightLDAP service.
resource "elestio_lightldap" "demo_lightldap" {
  project_id    = "2500"
  server_name   = "demo-lightldap"
  version       = "stable"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
