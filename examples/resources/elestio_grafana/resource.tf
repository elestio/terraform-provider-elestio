# Create and manage Grafana service.
resource "elestio_grafana" "demo_grafana" {
  project_id    = "2500"
  server_name   = "demo-grafana"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
