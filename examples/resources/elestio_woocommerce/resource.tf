# Create and manage WooCommerce service.
resource "elestio_woocommerce" "demo_woocommerce" {
  project_id    = "2500"
  server_name   = "demo-woocommerce"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
