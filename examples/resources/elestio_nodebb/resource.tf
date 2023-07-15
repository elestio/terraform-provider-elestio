# Create and manage NodeBB service.
resource "elestio_nodebb" "demo_nodebb" {
  project_id    = "2500"
  server_name   = "demo-nodebb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
