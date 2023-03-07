# Create and manage Zulip service.
resource "elestio_zulip" "my_zulip" {
  project_id    = "2500"
  server_name   = "awesome-zulip"
  server_type   = "SMALL-1C-2G"
  version       = "5.1-0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
