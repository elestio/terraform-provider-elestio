# Create and manage Couchbase service.
resource "elestio_couchbase" "demo_couchbase" {
  project_id    = "2500"
  server_name   = "demo-couchbase"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
