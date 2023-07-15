# Create and manage Keycloak service.
resource "elestio_keycloak" "demo_keycloak" {
  project_id    = "2500"
  server_name   = "demo-keycloak"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
