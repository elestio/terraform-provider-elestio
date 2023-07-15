# Create and manage Appsmith service.
resource "elestio_appsmith" "demo_appsmith" {
  project_id    = "2500"
  server_name   = "demo-appsmith"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
