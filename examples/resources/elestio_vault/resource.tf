# Create and manage Vault service.
resource "elestio_vault" "demo_vault" {
  project_id    = "2500"
  server_name   = "demo-vault"
  version       = "1.13.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
