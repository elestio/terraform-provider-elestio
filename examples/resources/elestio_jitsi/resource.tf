# Create and manage Jitsi service.
resource "elestio_jitsi" "my_jitsi" {
  project_id    = "2500"
  server_name   = "awesome-jitsi"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
