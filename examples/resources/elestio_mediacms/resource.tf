# Create and manage MediaCMS service.
resource "elestio_mediacms" "demo_mediacms" {
  project_id    = "2500"
  server_name   = "demo-mediacms"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
