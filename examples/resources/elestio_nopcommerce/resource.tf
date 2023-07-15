# Create and manage nopCommerce service.
resource "elestio_nopcommerce" "demo_nopcommerce" {
  project_id    = "2500"
  server_name   = "demo-nopcommerce"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
