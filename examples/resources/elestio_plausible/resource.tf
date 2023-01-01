# Create and manage Plausible Analytics service.
resource "elestio_plausible" "my_plausible" {
  project_id    = "2500"
  server_name   = "awesome-plausible"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
