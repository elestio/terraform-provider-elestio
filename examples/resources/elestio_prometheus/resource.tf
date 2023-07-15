# Create and manage Prometheus service.
resource "elestio_prometheus" "demo_prometheus" {
  project_id    = "2500"
  server_name   = "demo-prometheus"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
