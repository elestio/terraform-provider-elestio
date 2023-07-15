# Create and manage Strapi service.
resource "elestio_strapi" "demo_strapi" {
  project_id    = "2500"
  server_name   = "demo-strapi"
  version       = "3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
