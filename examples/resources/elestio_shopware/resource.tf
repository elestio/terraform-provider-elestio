# Create and manage Shopware service.
resource "elestio_shopware" "demo_shopware" {
  project_id    = "2500"
  server_name   = "demo-shopware"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
