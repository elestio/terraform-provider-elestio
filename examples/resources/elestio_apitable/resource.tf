# Create and manage APITable service.
resource "elestio_apitable" "demo_apitable" {
  project_id    = "2500"
  server_name   = "demo-apitable"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
