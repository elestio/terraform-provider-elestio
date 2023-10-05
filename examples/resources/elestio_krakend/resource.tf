# Create and manage KrakenD service.
resource "elestio_krakend" "demo_krakend" {
  project_id    = "2500"
  server_name   = "demo-krakend"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
