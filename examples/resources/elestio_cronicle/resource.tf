# Create and manage Cronicle service.
resource "elestio_cronicle" "demo_cronicle" {
  project_id    = "2500"
  server_name   = "demo-cronicle"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
