# Create and manage Flarum service.
resource "elestio_flarum" "demo_flarum" {
  project_id    = "2500"
  server_name   = "demo-flarum"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
