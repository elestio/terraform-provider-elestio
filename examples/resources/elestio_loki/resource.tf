# Create and manage Loki service.
resource "elestio_loki" "demo_loki" {
  project_id    = "2500"
  server_name   = "demo-loki"
  version       = "2.8.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
