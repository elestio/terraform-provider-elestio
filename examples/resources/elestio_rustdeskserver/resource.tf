# Create and manage RustdeskServer service.
resource "elestio_rustdeskserver" "demo_rustdeskserver" {
  project_id    = "2500"
  server_name   = "demo-rustdeskserver"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
