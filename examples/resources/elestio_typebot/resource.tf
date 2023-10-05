# Create and manage Typebot service.
resource "elestio_typebot" "demo_typebot" {
  project_id    = "2500"
  server_name   = "demo-typebot"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
