# Create and manage Ackee service.
resource "elestio_ackee" "demo_ackee" {
  project_id    = "2500"
  server_name   = "demo-ackee"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
