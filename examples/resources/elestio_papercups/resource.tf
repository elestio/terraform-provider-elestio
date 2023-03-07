# Create and manage Papercups service.
resource "elestio_papercups" "my_papercups" {
  project_id    = "2500"
  server_name   = "awesome-papercups"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
