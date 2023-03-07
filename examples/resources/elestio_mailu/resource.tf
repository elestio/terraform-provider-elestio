# Create and manage Mailu service.
resource "elestio_mailu" "my_mailu" {
  project_id    = "2500"
  server_name   = "awesome-mailu"
  server_type   = "SMALL-1C-2G"
  version       = "1.9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
