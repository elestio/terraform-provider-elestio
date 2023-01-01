# Create and manage Strapi service.
resource "elestio_strapi" "my_strapi" {
  project_id    = "2500"
  server_name   = "awesome-strapi"
  server_type   = "SMALL-1C-2G"
  version       = "3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
