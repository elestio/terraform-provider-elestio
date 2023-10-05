# Create and manage Bigcapital service.
resource "elestio_bigcapital" "demo_bigcapital" {
  project_id    = "2500"
  server_name   = "demo-bigcapital"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
