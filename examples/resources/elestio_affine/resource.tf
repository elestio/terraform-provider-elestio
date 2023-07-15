# Create and manage Affine service.
resource "elestio_affine" "demo_affine" {
  project_id    = "2500"
  server_name   = "demo-affine"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
