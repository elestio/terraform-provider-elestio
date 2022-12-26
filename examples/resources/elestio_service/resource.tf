# Create and manage a service.
resource "elestio_service" "service1" {
  project_id    = "project_id"
  server_name   = "service1"
  server_type   = "SMALL-1C-2G"
  template_id   = 11 // postgres
  version       = "14"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "admin@exemple.com"
}
