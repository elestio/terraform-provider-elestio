# Create and manage Fugu service.
resource "elestio_fugu" "demo_fugu" {
  project_id    = "2500"
  server_name   = "demo-fugu"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
