# Create and manage ToolJet service.
resource "elestio_tooljet" "my_tooljet" {
  project_id    = "2500"
  server_name   = "awesome-tooljet"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
