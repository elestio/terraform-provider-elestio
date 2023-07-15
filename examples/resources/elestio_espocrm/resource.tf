# Create and manage EspoCRM service.
resource "elestio_espocrm" "demo_espocrm" {
  project_id    = "2500"
  server_name   = "demo-espocrm"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
