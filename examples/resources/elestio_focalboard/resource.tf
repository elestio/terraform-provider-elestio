# Create and manage FocalBoard service.
resource "elestio_focalboard" "my_focalboard" {
  project_id    = "2500"
  server_name   = "awesome-focalboard"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
