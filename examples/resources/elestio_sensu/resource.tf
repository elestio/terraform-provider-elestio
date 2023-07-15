# Create and manage Sensu service.
resource "elestio_sensu" "demo_sensu" {
  project_id    = "2500"
  server_name   = "demo-sensu"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
