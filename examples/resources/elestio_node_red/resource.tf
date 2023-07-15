# Create and manage Node-red service.
resource "elestio_node_red" "demo_node_red" {
  project_id    = "2500"
  server_name   = "demo-node_red"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
