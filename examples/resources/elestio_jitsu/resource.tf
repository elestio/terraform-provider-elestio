# Create and manage Jitsu service.
resource "elestio_jitsu" "demo_jitsu" {
  project_id    = "2500"
  server_name   = "demo-jitsu"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
