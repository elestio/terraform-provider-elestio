# Create and manage Daptin service.
resource "elestio_daptin" "demo_daptin" {
  project_id    = "2500"
  server_name   = "demo-daptin"
  version       = "master"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
