# Create and manage FreeScout service.
resource "elestio_freescout" "demo_freescout" {
  project_id    = "2500"
  server_name   = "demo-freescout"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
