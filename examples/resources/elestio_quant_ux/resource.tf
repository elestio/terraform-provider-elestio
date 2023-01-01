# Create and manage Quant-UX service.
resource "elestio_quant_ux" "my_quant_ux" {
  project_id    = "2500"
  server_name   = "awesome-quant_ux"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
