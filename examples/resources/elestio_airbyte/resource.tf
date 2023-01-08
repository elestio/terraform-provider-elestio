# Create and manage Airbyte service.
resource "elestio_airbyte" "my_airbyte" {
  project_id    = "2500"
  server_name   = "awesome-airbyte"
  server_type   = "SMALL-1C-2G"
  version       = "0.40.17"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
