# Create and manage Directus service.
resource "elestio_directus" "my_directus" {
  project_id    = "2500"
  server_name   = "awesome-directus"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
