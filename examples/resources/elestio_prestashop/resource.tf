# Create and manage Prestashop service.
resource "elestio_prestashop" "my_prestashop" {
  project_id    = "2500"
  server_name   = "awesome-prestashop"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
