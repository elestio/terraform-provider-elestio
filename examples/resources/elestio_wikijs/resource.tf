# Create and manage Wikijs service.
resource "elestio_wikijs" "demo_wikijs" {
  project_id    = "2500"
  server_name   = "demo-wikijs"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
