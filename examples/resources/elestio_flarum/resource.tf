# Create and manage Flarum service.
resource "elestio_flarum" "my_flarum" {
  project_id    = "2500"
  server_name   = "awesome-flarum"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
