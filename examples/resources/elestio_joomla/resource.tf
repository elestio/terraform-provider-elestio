# Create and manage Joomla service.
resource "elestio_joomla" "demo_joomla" {
  project_id    = "2500"
  server_name   = "demo-joomla"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
