# Create and manage Wordpress service.
resource "elestio_wordpress" "demo_wordpress" {
  project_id    = "2500"
  server_name   = "demo-wordpress"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
