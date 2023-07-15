# Create and manage NocoDB service.
resource "elestio_nocodb" "demo_nocodb" {
  project_id    = "2500"
  server_name   = "demo-nocodb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
