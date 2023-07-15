# Create and manage ToolJet service.
resource "elestio_tooljet" "demo_tooljet" {
  project_id    = "2500"
  server_name   = "demo-tooljet"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
