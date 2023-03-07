# Create and manage Metabase service.
resource "elestio_metabase" "my_metabase" {
  project_id    = "2500"
  server_name   = "awesome-metabase"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
