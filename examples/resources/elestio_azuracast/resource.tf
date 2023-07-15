# Create and manage AzuraCast service.
resource "elestio_azuracast" "demo_azuracast" {
  project_id    = "2500"
  server_name   = "demo-azuracast"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
