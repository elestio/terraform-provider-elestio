# Create and manage Onedev service.
resource "elestio_onedev" "my_onedev" {
  project_id    = "2500"
  server_name   = "awesome-onedev"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
