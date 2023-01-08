# Create and manage M3DB service.
resource "elestio_m3db" "my_m3db" {
  project_id    = "2500"
  server_name   = "awesome-m3db"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
