resource "elestio_meilisearch" "example" {
  project_id    = "2500"
  version       = "v1.14.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
