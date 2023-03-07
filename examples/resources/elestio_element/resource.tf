# Create and manage Element service.
resource "elestio_element" "my_element" {
  project_id    = "2500"
  server_name   = "awesome-element"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
