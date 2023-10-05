# Create and manage Lobsters service.
resource "elestio_lobsters" "demo_lobsters" {
  project_id    = "2500"
  server_name   = "demo-lobsters"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
