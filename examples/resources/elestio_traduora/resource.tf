# Create and manage Traduora service.
resource "elestio_traduora" "my_traduora" {
  project_id    = "2500"
  server_name   = "awesome-traduora"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
