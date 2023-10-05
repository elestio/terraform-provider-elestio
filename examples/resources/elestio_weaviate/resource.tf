# Create and manage Weaviate service.
resource "elestio_weaviate" "demo_weaviate" {
  project_id    = "2500"
  server_name   = "demo-weaviate"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
