# Create and manage Mage AI service.
resource "elestio_mage" "demo_mage" {
  project_id    = "2500"
  server_name   = "demo-mage"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
