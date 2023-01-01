# Create and manage Vault service.
resource "elestio_vault" "my_vault" {
  project_id    = "2500"
  server_name   = "awesome-vault"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
