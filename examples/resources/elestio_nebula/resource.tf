# Create and manage Nebula service.
resource "elestio_nebula" "demo_nebula" {
  project_id    = "2500"
  server_name   = "demo-nebula"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
