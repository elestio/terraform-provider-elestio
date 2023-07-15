# Create and manage OpenSearch service.
resource "elestio_opensearch" "demo_opensearch" {
  project_id    = "2500"
  server_name   = "demo-opensearch"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
