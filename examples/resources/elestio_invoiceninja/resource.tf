# Create and manage InvoiceNinja service.
resource "elestio_invoiceninja" "my_invoiceninja" {
  project_id    = "2500"
  server_name   = "awesome-invoiceninja"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
