# Create and manage Jupyter Notebook service.
resource "elestio_jupyter" "demo_jupyter" {
  project_id    = "2500"
  server_name   = "demo-jupyter"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
