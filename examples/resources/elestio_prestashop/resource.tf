# Create and manage Prestashop service.
resource "elestio_prestashop" "demo_prestashop" {
  project_id    = "2500"
  server_name   = "demo-prestashop"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
