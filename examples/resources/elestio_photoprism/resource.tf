# Create and manage PhotoPrism service.
resource "elestio_photoprism" "demo_photoprism" {
  project_id    = "2500"
  server_name   = "demo-photoprism"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
