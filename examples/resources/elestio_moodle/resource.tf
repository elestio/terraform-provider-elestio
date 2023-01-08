# Create and manage Moodle service.
resource "elestio_moodle" "my_moodle" {
  project_id    = "2500"
  server_name   = "awesome-moodle"
  server_type   = "SMALL-1C-2G"
  version       = "3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
