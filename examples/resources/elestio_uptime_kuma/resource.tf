# Create and manage Uptime-kuma service.
resource "elestio_uptime_kuma" "my_uptime_kuma" {
  project_id    = "2500"
  server_name   = "awesome-uptime_kuma"
  server_type   = "SMALL-1C-2G"
  version       = "1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
