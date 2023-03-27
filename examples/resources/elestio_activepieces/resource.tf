# Create and manage Activepieces service.
resource "elestio_activepieces" "my_activepieces" {
  project_id    = "2500"
  server_name   = "awesome-activepieces"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
