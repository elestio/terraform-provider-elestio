# Create and manage Mautic service.
resource "elestio_mautic" "demo_mautic" {
  project_id    = "2500"
  server_name   = "demo-mautic"
  version       = "v4"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
