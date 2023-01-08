# Create and manage Browserless service.
resource "elestio_browserless" "my_browserless" {
  project_id    = "2500"
  server_name   = "awesome-browserless"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
