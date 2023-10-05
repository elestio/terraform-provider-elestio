# Create and manage ILLA service.
resource "elestio_illa" "demo_illa" {
  project_id    = "2500"
  server_name   = "demo-illa"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
