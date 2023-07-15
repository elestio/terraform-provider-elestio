# Create and manage Element service.
resource "elestio_element" "demo_element" {
  project_id    = "2500"
  server_name   = "demo-element"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
