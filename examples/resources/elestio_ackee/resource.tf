# Create and manage Ackee service.
resource "elestio_ackee" "my_ackee" {
  project_id    = "2500"
  server_name   = "awesome-ackee"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
