# Create and manage Nagios service.
resource "elestio_nagios" "demo_nagios" {
  project_id    = "2500"
  server_name   = "demo-nagios"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
