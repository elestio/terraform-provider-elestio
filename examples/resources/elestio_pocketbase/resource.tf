# Create and manage PocketBase service.
resource "elestio_pocketbase" "demo_pocketbase" {
  project_id    = "2500"
  server_name   = "demo-pocketbase"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
