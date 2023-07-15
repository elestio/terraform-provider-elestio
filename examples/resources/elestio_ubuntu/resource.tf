# Create and manage Ubuntu service.
resource "elestio_ubuntu" "demo_ubuntu" {
  project_id    = "2500"
  server_name   = "demo-ubuntu"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
