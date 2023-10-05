# Create and manage Bytebase service.
resource "elestio_bytebase" "demo_bytebase" {
  project_id    = "2500"
  server_name   = "demo-bytebase"
  version       = "1.17.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
