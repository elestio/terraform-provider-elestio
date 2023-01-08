# Create and manage Keycloak service.
resource "elestio_keycloak" "my_keycloak" {
  project_id    = "2500"
  server_name   = "awesome-keycloak"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
