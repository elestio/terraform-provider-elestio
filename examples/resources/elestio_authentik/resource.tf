# Create and manage Authentik service.
resource "elestio_authentik" "demo_authentik" {
  project_id    = "2500"
  server_name   = "demo-authentik"
  version       = "2023.3.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
