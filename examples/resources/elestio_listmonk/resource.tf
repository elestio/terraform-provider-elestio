# Create and manage Listmonk service.
resource "elestio_listmonk" "demo_listmonk" {
  project_id    = "2500"
  server_name   = "demo-listmonk"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
