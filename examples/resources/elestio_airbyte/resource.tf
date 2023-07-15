# Create and manage Airbyte service.
resource "elestio_airbyte" "demo_airbyte" {
  project_id    = "2500"
  server_name   = "demo-airbyte"
  version       = "0.43.1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
