# Create and manage Traggo service.
resource "elestio_traggo" "demo_traggo" {
  project_id    = "2500"
  server_name   = "demo-traggo"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
