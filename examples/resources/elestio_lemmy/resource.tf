# Create and manage Lemmy service.
resource "elestio_lemmy" "demo_lemmy" {
  project_id    = "2500"
  server_name   = "demo-lemmy"
  version       = "0.17.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
