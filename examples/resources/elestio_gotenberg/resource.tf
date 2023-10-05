# Create and manage Gotenberg service.
resource "elestio_gotenberg" "demo_gotenberg" {
  project_id    = "2500"
  server_name   = "demo-gotenberg"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
