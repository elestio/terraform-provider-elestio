# Create and manage Quant-UX service.
resource "elestio_quant_ux" "demo_quant_ux" {
  project_id    = "2500"
  server_name   = "demo-quant_ux"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
