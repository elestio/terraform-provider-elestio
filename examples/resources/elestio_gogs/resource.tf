# Create and manage Gogs service.
resource "elestio_gogs" "my_gogs" {
  project_id    = "2500"
  server_name   = "awesome-gogs"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
