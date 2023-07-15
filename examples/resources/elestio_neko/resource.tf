# Create and manage Neko Rooms service.
resource "elestio_neko" "demo_neko" {
  project_id    = "2500"
  server_name   = "demo-neko"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
