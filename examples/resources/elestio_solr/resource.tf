# Create and manage Solr service.
resource "elestio_solr" "demo_solr" {
  project_id    = "2500"
  server_name   = "demo-solr"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
