# Create and manage MetaTrader5 service.
resource "elestio_metatrader5" "my_metatrader5" {
  project_id    = "2500"
  server_name   = "awesome-metatrader5"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
