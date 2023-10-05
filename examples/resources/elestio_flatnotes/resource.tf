# Create and manage Flatnotes service.
resource "elestio_flatnotes" "demo_flatnotes" {
  project_id    = "2500"
  server_name   = "demo-flatnotes"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
