# Create and manage DockerRegistry service.
resource "elestio_dockerregistry" "demo_dockerregistry" {
  project_id    = "2500"
  server_name   = "demo-dockerregistry"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
