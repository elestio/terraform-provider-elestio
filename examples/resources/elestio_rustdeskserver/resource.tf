# Create and manage RustdeskServer service.
resource "elestio_rustdeskserver" "my_rustdeskserver" {
  project_id    = "2500"
  server_name   = "awesome-rustdeskserver"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
