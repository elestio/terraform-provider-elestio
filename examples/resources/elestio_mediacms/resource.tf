# Create and manage MediaCMS service.
resource "elestio_mediacms" "my_mediacms" {
  project_id    = "2500"
  server_name   = "awesome-mediacms"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
