# Create and manage Parse service.
resource "elestio_parse" "my_parse" {
  project_id    = "2500"
  server_name   = "awesome-parse"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
