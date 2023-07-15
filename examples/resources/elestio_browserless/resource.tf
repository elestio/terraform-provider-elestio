# Create and manage Browserless service.
resource "elestio_browserless" "demo_browserless" {
  project_id    = "2500"
  server_name   = "demo-browserless"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
