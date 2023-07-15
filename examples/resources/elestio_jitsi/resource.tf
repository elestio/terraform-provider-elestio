# Create and manage Jitsi service.
resource "elestio_jitsi" "demo_jitsi" {
  project_id    = "2500"
  server_name   = "demo-jitsi"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
