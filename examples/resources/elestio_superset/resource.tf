# Create and manage Superset service.
resource "elestio_superset" "my_superset" {
  project_id    = "2500"
  server_name   = "awesome-superset"
  server_type   = "SMALL-1C-2G"
  version       = "master"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
