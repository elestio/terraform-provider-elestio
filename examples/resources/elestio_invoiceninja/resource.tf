# Create and manage InvoiceNinja service.
resource "elestio_invoiceninja" "demo_invoiceninja" {
  project_id    = "2500"
  server_name   = "demo-invoiceninja"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
