# Create and manage N8N service.
resource "elestio_n8n" "my_n8n" {
  project_id    = "2500"
  server_name   = "awesome-n8n"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
