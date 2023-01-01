# Create and manage Rocket.Chat service.
resource "elestio_rocket_chat" "my_rocket_chat" {
  project_id    = "2500"
  server_name   = "awesome-rocket_chat"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
