# Create and manage Infisical service.
resource "elestio_infisical" "demo_infisical" {
  project_id    = "2500"
  server_name   = "demo-infisical"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
