# Create and manage Castopod service.
resource "elestio_castopod" "demo_castopod" {
  project_id    = "2500"
  server_name   = "demo-castopod"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
