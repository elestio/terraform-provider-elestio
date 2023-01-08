# Create and manage Lightdash service.
resource "elestio_lightdash" "my_lightdash" {
  project_id    = "2500"
  server_name   = "awesome-lightdash"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
