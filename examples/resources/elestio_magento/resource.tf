# Create and manage Magento service.
resource "elestio_magento" "demo_magento" {
  project_id    = "2500"
  server_name   = "demo-magento"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
