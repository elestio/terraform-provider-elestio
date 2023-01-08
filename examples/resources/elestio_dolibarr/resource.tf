# Create and manage Dolibarr service.
resource "elestio_dolibarr" "my_dolibarr" {
  project_id    = "2500"
  server_name   = "awesome-dolibarr"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
