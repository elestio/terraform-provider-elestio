# Create and manage Ubuntu-Desktop service.
resource "elestio_ubuntu_desktop" "my_ubuntu_desktop" {
  project_id    = "2500"
  server_name   = "awesome-ubuntu_desktop"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
