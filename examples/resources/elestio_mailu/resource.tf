# Create and manage Mailu service.
resource "elestio_mailu" "demo_mailu" {
  project_id    = "2500"
  server_name   = "demo-mailu"
  version       = "2.0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
