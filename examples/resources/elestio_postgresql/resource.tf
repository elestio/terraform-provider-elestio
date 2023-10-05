# Create and manage PostgreSQL service.
resource "elestio_postgresql" "demo_postgresql" {
  project_id    = "2500"
  server_name   = "demo-postgresql"
  version       = "16"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
