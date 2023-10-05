# Create and manage Flagsmith service.
resource "elestio_flagsmith" "demo_flagsmith" {
  project_id    = "2500"
  server_name   = "demo-flagsmith"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
