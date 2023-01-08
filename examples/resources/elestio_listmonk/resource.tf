# Create and manage Listmonk service.
resource "elestio_listmonk" "my_listmonk" {
  project_id    = "2500"
  server_name   = "awesome-listmonk"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
