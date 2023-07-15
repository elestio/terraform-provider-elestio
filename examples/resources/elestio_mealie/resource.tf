# Create and manage Mealie service.
resource "elestio_mealie" "demo_mealie" {
  project_id    = "2500"
  server_name   = "demo-mealie"
  version       = "omni-nightly"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
