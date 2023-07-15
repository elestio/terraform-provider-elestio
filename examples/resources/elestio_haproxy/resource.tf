# Create and manage HAProxy service.
resource "elestio_haproxy" "demo_haproxy" {
  project_id    = "2500"
  server_name   = "demo-haproxy"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
