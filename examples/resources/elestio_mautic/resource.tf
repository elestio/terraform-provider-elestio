# Create and manage Mautic service.
resource "elestio_mautic" "my_mautic" {
  project_id    = "2500"
  server_name   = "awesome-mautic"
  server_type   = "SMALL-1C-2G"
  version       = "v4"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
