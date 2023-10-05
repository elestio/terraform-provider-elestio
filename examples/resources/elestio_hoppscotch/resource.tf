# Create and manage Hoppscotch service.
resource "elestio_hoppscotch" "demo_hoppscotch" {
  project_id    = "2500"
  server_name   = "demo-hoppscotch"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
