# Create and manage Tracardi service.
resource "elestio_tracardi" "demo_tracardi" {
  project_id    = "2500"
  server_name   = "demo-tracardi"
  version       = "0.7.1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
