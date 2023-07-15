# Create and manage Ghost service.
resource "elestio_ghost" "demo_ghost" {
  project_id    = "2500"
  server_name   = "demo-ghost"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
