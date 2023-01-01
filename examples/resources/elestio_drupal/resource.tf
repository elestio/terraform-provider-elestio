# Create and manage Drupal service.
resource "elestio_drupal" "my_drupal" {
  project_id    = "2500"
  server_name   = "awesome-drupal"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
