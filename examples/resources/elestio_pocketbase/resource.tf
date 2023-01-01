# Create and manage PocketBase service.
resource "elestio_pocketbase" "my_pocketbase" {
  project_id    = "2500"
  server_name   = "awesome-pocketbase"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
