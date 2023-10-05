# Create and manage Uptime-kuma service.
resource "elestio_uptime_kuma" "demo_uptime_kuma" {
  project_id    = "2500"
  server_name   = "demo-uptime_kuma"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
