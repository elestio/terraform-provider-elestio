# Create and manage SuiteCRM service.
resource "elestio_suitecrm" "my_suitecrm" {
  project_id    = "2500"
  server_name   = "awesome-suitecrm"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
