# Create and manage Dittofeed service.
resource "elestio_dittofeed" "demo_dittofeed" {
  project_id    = "2500"
  server_name   = "demo-dittofeed"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
