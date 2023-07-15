# Create and manage Odoo ERP [TEMPLATE_DOCUMENTATION_NAME] CRM service.
resource "elestio_odoo" "demo_odoo" {
  project_id    = "2500"
  server_name   = "demo-odoo"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
