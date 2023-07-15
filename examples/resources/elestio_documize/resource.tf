# Create and manage Documize service.
resource "elestio_documize" "demo_documize" {
  project_id    = "2500"
  server_name   = "demo-documize"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
