# Create and manage Imagor service.
resource "elestio_imagor" "demo_imagor" {
  project_id    = "2500"
  server_name   = "demo-imagor"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
