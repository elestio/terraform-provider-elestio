# Create and manage Drupal service.
resource "elestio_drupal" "demo_drupal" {
  project_id    = "2500"
  server_name   = "demo-drupal"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
