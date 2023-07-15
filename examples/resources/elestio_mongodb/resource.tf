# Create and manage MongoDB service.
resource "elestio_mongodb" "demo_mongodb" {
  project_id    = "2500"
  server_name   = "demo-mongodb"
  version       = "6"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
