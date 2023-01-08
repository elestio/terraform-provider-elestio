# Create and manage Mastodon service.
resource "elestio_mastodon" "my_mastodon" {
  project_id    = "2500"
  server_name   = "awesome-mastodon"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
