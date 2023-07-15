# Create and manage Dolibarr service.
resource "elestio_dolibarr" "demo_dolibarr" {
  project_id    = "2500"
  server_name   = "demo-dolibarr"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
