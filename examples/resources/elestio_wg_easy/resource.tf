# Create and manage WG-Easy service.
resource "elestio_wg_easy" "demo_wg_easy" {
  project_id    = "2500"
  server_name   = "demo-wg_easy"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
