# Create and manage OpenSlides service.
resource "elestio_openslides" "demo_openslides" {
  project_id    = "2500"
  server_name   = "demo-openslides"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
