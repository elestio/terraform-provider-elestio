# Create and manage Mattermost Team Edition service.
resource "elestio_mattermost" "demo_mattermost" {
  project_id    = "2500"
  server_name   = "demo-mattermost"
  version       = "master"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
