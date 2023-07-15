# Create and manage Knime service.
resource "elestio_knime" "demo_knime" {
  project_id    = "2500"
  server_name   = "demo-knime"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
