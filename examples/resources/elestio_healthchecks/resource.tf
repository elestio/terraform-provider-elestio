# Create and manage Healthchecks service.
resource "elestio_healthchecks" "demo_healthchecks" {
  project_id    = "2500"
  server_name   = "demo-healthchecks"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
