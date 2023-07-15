# Create and manage Gerrit service.
resource "elestio_gerrit" "demo_gerrit" {
  project_id    = "2500"
  server_name   = "demo-gerrit"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
