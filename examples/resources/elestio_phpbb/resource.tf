# Create and manage PhpBB service.
resource "elestio_phpbb" "demo_phpbb" {
  project_id    = "2500"
  server_name   = "demo-phpbb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
