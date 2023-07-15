# Create and manage OpenResty service.
resource "elestio_openresty" "demo_openresty" {
  project_id    = "2500"
  server_name   = "demo-openresty"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
