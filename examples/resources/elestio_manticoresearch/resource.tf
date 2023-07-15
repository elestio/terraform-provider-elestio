# Create and manage ManticoreSearch service.
resource "elestio_manticoresearch" "demo_manticoresearch" {
  project_id    = "2500"
  server_name   = "demo-manticoresearch"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
