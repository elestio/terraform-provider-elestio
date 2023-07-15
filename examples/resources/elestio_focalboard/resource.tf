# Create and manage FocalBoard service.
resource "elestio_focalboard" "demo_focalboard" {
  project_id    = "2500"
  server_name   = "demo-focalboard"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
