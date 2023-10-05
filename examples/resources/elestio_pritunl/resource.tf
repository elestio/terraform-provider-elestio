# Create and manage Pritunl service.
resource "elestio_pritunl" "demo_pritunl" {
  project_id    = "2500"
  server_name   = "demo-pritunl"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
