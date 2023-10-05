# Create and manage Friendica service.
resource "elestio_friendica" "demo_friendica" {
  project_id    = "2500"
  server_name   = "demo-friendica"
  version       = "stable"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
