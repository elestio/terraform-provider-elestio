# Create and manage Linux-desktop service.
resource "elestio_linux_desktop" "my_linux_desktop" {
  project_id    = "2500"
  server_name   = "awesome-linux_desktop"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
