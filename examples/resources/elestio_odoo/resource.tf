# Create and manage Odoo ERP [TEMPLATE_DOCUMENTATION_NAME] CRM service.
resource "elestio_odoo" "my_odoo" {
  project_id    = "2500"
  server_name   = "awesome-odoo"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
