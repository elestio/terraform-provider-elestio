# Create and manage Neo4j service.
resource "elestio_neo4j" "my_neo4j" {
  project_id    = "2500"
  server_name   = "awesome-neo4j"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
