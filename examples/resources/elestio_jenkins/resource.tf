# Create and manage Jenkins service.
resource "elestio_jenkins" "demo_jenkins" {
  project_id    = "2500"
  server_name   = "demo-jenkins"
  version       = "lts-jdk11"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
