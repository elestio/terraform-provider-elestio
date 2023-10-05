# Create and manage ZincSearch service.
resource "elestio_zincsearch" "demo_zincsearch" {
  project_id    = "2500"
  server_name   = "demo-zincsearch"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
