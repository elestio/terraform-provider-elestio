# Create and manage Chaskiq service.
resource "elestio_chaskiq" "demo_chaskiq" {
  project_id    = "2500"
  server_name   = "demo-chaskiq"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
