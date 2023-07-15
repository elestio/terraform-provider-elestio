# Create and manage Nomad service.
resource "elestio_nomad" "demo_nomad" {
  project_id    = "2500"
  server_name   = "demo-nomad"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
