# Create and manage M3DB service.
resource "elestio_m3db" "demo_m3db" {
  project_id    = "2500"
  server_name   = "demo-m3db"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
