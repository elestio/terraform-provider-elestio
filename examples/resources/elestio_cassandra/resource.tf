# Create and manage Cassandra service.
resource "elestio_cassandra" "demo_cassandra" {
  project_id    = "2500"
  server_name   = "demo-cassandra"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
