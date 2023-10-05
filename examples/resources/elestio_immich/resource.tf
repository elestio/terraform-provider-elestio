# Create and manage Immich service.
resource "elestio_immich" "demo_immich" {
  project_id    = "2500"
  server_name   = "demo-immich"
  version       = "release"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
