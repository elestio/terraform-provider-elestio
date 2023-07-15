# Create and manage Traduora service.
resource "elestio_traduora" "demo_traduora" {
  project_id    = "2500"
  server_name   = "demo-traduora"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
