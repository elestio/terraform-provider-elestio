# Create and manage ScyllaDB service.
resource "elestio_scylladb" "demo_scylladb" {
  project_id    = "2500"
  server_name   = "demo-scylladb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
