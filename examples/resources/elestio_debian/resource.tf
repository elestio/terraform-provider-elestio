# Create and manage Debian service.
resource "elestio_debian" "demo_debian" {
  project_id    = "2500"
  server_name   = "demo-debian"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
