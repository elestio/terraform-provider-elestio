# Create and manage Botpress service.
resource "elestio_botpress" "demo_botpress" {
  project_id    = "2500"
  server_name   = "demo-botpress"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
