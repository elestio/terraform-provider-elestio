# Create and manage Canopsis service.
resource "elestio_canopsis" "my_canopsis" {
  project_id    = "2500"
  server_name   = "awesome-canopsis"
  server_type   = "SMALL-1C-2G"
  version       = "4.3.9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
