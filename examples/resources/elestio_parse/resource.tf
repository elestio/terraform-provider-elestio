# Create and manage Parse service.
resource "elestio_parse" "demo_parse" {
  project_id    = "2500"
  server_name   = "demo-parse"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
