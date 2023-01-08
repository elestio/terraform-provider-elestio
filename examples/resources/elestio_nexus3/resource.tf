# Create and manage Nexus3 service.
resource "elestio_nexus3" "my_nexus3" {
  project_id    = "2500"
  server_name   = "awesome-nexus3"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
