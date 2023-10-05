# Create and manage Hop service.
resource "elestio_hop" "demo_hop" {
  project_id    = "2500"
  server_name   = "demo-hop"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
