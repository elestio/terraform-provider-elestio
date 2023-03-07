# Create and manage CryptPad service.
resource "elestio_cryptpad" "my_cryptpad" {
  project_id    = "2500"
  server_name   = "awesome-cryptpad"
  server_type   = "SMALL-1C-2G"
  version       = "nginx"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
