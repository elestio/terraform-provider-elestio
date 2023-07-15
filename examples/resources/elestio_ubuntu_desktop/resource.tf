# Create and manage Ubuntu-Desktop service.
resource "elestio_ubuntu_desktop" "demo_ubuntu_desktop" {
  project_id    = "2500"
  server_name   = "demo-ubuntu_desktop"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
