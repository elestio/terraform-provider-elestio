# Create and manage CryptPad service.
resource "elestio_cryptpad" "demo_cryptpad" {
  project_id    = "2500"
  server_name   = "demo-cryptpad"
  version       = "nginx"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
