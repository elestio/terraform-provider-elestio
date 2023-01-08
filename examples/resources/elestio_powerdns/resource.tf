# Create and manage PowerDNS service.
resource "elestio_powerdns" "my_powerdns" {
  project_id    = "2500"
  server_name   = "awesome-powerdns"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
