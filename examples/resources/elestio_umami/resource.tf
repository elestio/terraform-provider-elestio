# Create and manage Umami service.
resource "elestio_umami" "demo_umami" {
  project_id    = "2500"
  server_name   = "demo-umami"
  version       = "postgresql-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
