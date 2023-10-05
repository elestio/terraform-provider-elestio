# Create and manage Hydra service.
resource "elestio_hydra" "demo_hydra" {
  project_id    = "2500"
  server_name   = "demo-hydra"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
