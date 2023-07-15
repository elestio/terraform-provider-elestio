# Create and manage Vaultwarden service.
resource "elestio_vaultwarden" "demo_vaultwarden" {
  project_id    = "2500"
  server_name   = "demo-vaultwarden"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
