# Create and manage Papercups service.
resource "elestio_papercups" "demo_papercups" {
  project_id    = "2500"
  server_name   = "demo-papercups"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
