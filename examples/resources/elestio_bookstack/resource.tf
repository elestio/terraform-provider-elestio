# Create and manage BookStack service.
resource "elestio_bookstack" "demo_bookstack" {
  project_id    = "2500"
  server_name   = "demo-bookstack"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
