# Create and manage PhotoPrism service.
resource "elestio_photoprism" "my_photoprism" {
  project_id    = "2500"
  server_name   = "awesome-photoprism"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
