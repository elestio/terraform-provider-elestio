# Create and manage Appwrite service.
resource "elestio_appwrite" "my_appwrite" {
  project_id    = "2500"
  server_name   = "awesome-appwrite"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
