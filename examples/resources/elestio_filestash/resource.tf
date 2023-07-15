# Create and manage Filestash service.
resource "elestio_filestash" "demo_filestash" {
  project_id    = "2500"
  server_name   = "demo-filestash"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
