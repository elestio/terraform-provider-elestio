# Create and manage Rocket.Chat service.
resource "elestio_rocket_chat" "demo_rocket_chat" {
  project_id    = "2500"
  server_name   = "demo-rocket_chat"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
