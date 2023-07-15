# Create and manage Zammad service.
resource "elestio_zammad" "demo_zammad" {
  project_id    = "2500"
  server_name   = "demo-zammad"
  version       = "5.4.1-38"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
