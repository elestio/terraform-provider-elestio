# Create and manage Moodle service.
resource "elestio_moodle" "demo_moodle" {
  project_id    = "2500"
  server_name   = "demo-moodle"
  version       = "4.2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
