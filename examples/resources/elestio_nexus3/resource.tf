# Create and manage Nexus3 service.
resource "elestio_nexus3" "demo_nexus3" {
  project_id    = "2500"
  server_name   = "demo-nexus3"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
