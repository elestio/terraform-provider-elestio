# Create and manage MetaTrader5 service.
resource "elestio_metatrader5" "demo_metatrader5" {
  project_id    = "2500"
  server_name   = "demo-metatrader5"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
