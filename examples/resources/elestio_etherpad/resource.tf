# Create and manage Etherpad service.
resource "elestio_etherpad" "demo_etherpad" {
  project_id    = "2500"
  server_name   = "demo-etherpad"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
