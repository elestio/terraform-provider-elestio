# Create and manage Picoshare service.
resource "elestio_picoshare" "demo_picoshare" {
  project_id    = "2500"
  server_name   = "demo-picoshare"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
