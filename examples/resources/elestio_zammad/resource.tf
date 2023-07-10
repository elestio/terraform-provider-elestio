# Create and manage Zammad service.
resource "elestio_zammad" "my_zammad" {
  project_id    = "2500"
  server_name   = "awesome-zammad"
  server_type   = "SMALL-1C-2G"
  version       = "5.4.1-38"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
