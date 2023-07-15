# Create and manage MeiliSearch service.
resource "elestio_meilisearch" "demo_meilisearch" {
  project_id    = "2500"
  server_name   = "demo-meilisearch"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
