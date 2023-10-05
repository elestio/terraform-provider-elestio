# Create and manage Zitadel service.
resource "elestio_zitadel" "demo_zitadel" {
  project_id    = "2500"
  server_name   = "demo-zitadel"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
