# Create and manage MeiliSearch service.
resource "elestio_meilisearch" "my_meilisearch" {
  project_id    = "2500"
  server_name   = "awesome-meilisearch"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
