# Create and manage Chatwoot service.
resource "elestio_chatwoot" "demo_chatwoot" {
  project_id    = "2500"
  server_name   = "demo-chatwoot"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
