# Create and manage Authelia service.
resource "elestio_authelia" "demo_authelia" {
  project_id    = "2500"
  server_name   = "demo-authelia"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
