# Create and manage Typesense service.
resource "elestio_typesense" "my_typesense" {
  project_id    = "2500"
  server_name   = "awesome-typesense"
  server_type   = "SMALL-1C-2G"
  version       = "0.23.0.rc66"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
