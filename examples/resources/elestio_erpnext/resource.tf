# Create and manage ErpNext service.
resource "elestio_erpnext" "my_erpnext" {
  project_id    = "2500"
  server_name   = "awesome-erpnext"
  server_type   = "SMALL-1C-2G"
  version       = "v12.25.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
