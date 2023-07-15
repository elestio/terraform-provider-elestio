# Create and manage SuiteCRM service.
resource "elestio_suitecrm" "demo_suitecrm" {
  project_id    = "2500"
  server_name   = "demo-suitecrm"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
