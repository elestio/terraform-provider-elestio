# Create and manage Mattermost Team Edition service.
resource "elestio_mattermost" "my_mattermost" {
  project_id    = "2500"
  server_name   = "awesome-mattermost"
  server_type   = "SMALL-1C-2G"
  version       = "master"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
