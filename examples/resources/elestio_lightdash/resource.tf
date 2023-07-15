# Create and manage Lightdash service.
resource "elestio_lightdash" "demo_lightdash" {
  project_id    = "2500"
  server_name   = "demo-lightdash"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
