# Create and manage Gogs service.
resource "elestio_gogs" "demo_gogs" {
  project_id    = "2500"
  server_name   = "demo-gogs"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
