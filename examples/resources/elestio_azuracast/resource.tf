# Create and manage AzuraCast service.
resource "elestio_azuracast" "my_azuracast" {
  project_id    = "2500"
  server_name   = "awesome-azuracast"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
