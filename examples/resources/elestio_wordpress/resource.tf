# Create and manage Wordpress service.
resource "elestio_wordpress" "my_wordpress" {
  project_id    = "2500"
  server_name   = "awesome-wordpress"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
