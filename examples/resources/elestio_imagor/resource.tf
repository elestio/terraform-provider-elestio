# Create and manage Imagor service.
resource "elestio_imagor" "my_imagor" {
  project_id    = "2500"
  server_name   = "awesome-imagor"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
