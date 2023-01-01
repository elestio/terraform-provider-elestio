# Create and manage Magento service.
resource "elestio_magento" "my_magento" {
  project_id    = "2500"
  server_name   = "awesome-magento"
  server_type   = "SMALL-1C-2G"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
