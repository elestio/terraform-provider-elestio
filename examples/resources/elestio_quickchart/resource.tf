# Create and manage QuickChart service.
resource "elestio_quickchart" "demo_quickchart" {
  project_id    = "2500"
  server_name   = "demo-quickchart"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
