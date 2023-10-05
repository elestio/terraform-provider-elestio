# Create and manage OpenBlocks service.
resource "elestio_openblocks" "demo_openblocks" {
  project_id    = "2500"
  server_name   = "demo-openblocks"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
