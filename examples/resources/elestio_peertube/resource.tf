# Create and manage PeerTube service.
resource "elestio_peertube" "demo_peertube" {
  project_id    = "2500"
  server_name   = "demo-peertube"
  version       = "production-bullseye"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
