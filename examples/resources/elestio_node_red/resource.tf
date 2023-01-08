# Create and manage Node-red service.
resource "elestio_node_red" "my_node_red" {
  project_id    = "2500"
  server_name   = "awesome-node_red"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
