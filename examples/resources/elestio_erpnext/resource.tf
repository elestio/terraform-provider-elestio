# Create and manage ErpNext service.
resource "elestio_erpnext" "demo_erpnext" {
  project_id    = "2500"
  server_name   = "demo-erpnext"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
