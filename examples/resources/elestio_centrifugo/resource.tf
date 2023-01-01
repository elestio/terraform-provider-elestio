# Create and manage Centrifugo service.
resource "elestio_centrifugo" "my_centrifugo" {
  project_id    = "2500"
  server_name   = "awesome-centrifugo"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
