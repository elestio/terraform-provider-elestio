# Create and manage Healthchecks service.
resource "elestio_healthchecks" "my_healthchecks" {
  project_id    = "2500"
  server_name   = "awesome-healthchecks"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
