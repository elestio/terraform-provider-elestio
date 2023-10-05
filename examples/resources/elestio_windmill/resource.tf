# Create and manage Windmill service.
resource "elestio_windmill" "demo_windmill" {
  project_id    = "2500"
  server_name   = "demo-windmill"
  version       = "main"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
