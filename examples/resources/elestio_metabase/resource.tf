# Create and manage Metabase service.
resource "elestio_metabase" "demo_metabase" {
  project_id    = "2500"
  server_name   = "demo-metabase"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
