# Create and manage OnlyOffice service.
resource "elestio_onlyoffice" "my_onlyoffice" {
  project_id    = "2500"
  server_name   = "awesome-onlyoffice"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
