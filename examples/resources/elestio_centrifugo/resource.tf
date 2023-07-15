# Create and manage Centrifugo service.
resource "elestio_centrifugo" "demo_centrifugo" {
  project_id    = "2500"
  server_name   = "demo-centrifugo"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
