# Create and manage Hasura service.
resource "elestio_hasura" "demo_hasura" {
  project_id    = "2500"
  server_name   = "demo-hasura"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
