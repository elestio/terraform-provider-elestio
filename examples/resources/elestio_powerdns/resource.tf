# Create and manage PowerDNS service.
resource "elestio_powerdns" "demo_powerdns" {
  project_id    = "2500"
  server_name   = "demo-powerdns"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
