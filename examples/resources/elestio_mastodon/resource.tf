# Create and manage Mastodon service.
resource "elestio_mastodon" "demo_mastodon" {
  project_id    = "2500"
  server_name   = "demo-mastodon"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
