# Create and manage Rallly service.
resource "elestio_rallly" "demo_rallly" {
  project_id    = "2500"
  server_name   = "demo-rallly"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
