# Create and manage VictoriaMetrics service.
resource "elestio_victoriametrics" "demo_victoriametrics" {
  project_id    = "2500"
  server_name   = "demo-victoriametrics"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
