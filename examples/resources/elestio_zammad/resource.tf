# Create and manage Zammad service.
resource "elestio_zammad" "my_zammad" {
  project_id    = "2500"
  server_name   = "awesome-zammad"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
