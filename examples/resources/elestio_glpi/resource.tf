# Create and manage GLPI service.
resource "elestio_glpi" "demo_glpi" {
  project_id    = "2500"
  server_name   = "demo-glpi"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
