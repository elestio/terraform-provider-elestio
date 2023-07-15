# Create and manage Guacamole service.
resource "elestio_guacamole" "demo_guacamole" {
  project_id    = "2500"
  server_name   = "demo-guacamole"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
