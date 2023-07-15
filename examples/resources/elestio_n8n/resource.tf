# Create and manage N8N service.
resource "elestio_n8n" "demo_n8n" {
  project_id    = "2500"
  server_name   = "demo-n8n"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
