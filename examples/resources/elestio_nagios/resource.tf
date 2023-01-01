# Create and manage Nagios service.
resource "elestio_nagios" "my_nagios" {
  project_id    = "2500"
  server_name   = "awesome-nagios"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
