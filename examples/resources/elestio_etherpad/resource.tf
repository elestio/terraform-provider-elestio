# Create and manage Etherpad service.
resource "elestio_etherpad" "my_etherpad" {
  project_id    = "2500"
  server_name   = "awesome-etherpad"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
