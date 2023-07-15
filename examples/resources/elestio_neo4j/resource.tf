# Create and manage Neo4j service.
resource "elestio_neo4j" "demo_neo4j" {
  project_id    = "2500"
  server_name   = "demo-neo4j"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
