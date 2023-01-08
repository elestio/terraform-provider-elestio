# Create and manage Wikijs service.
resource "elestio_wikijs" "my_wikijs" {
  project_id    = "2500"
  server_name   = "awesome-wikijs"
  server_type   = "SMALL-1C-2G"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
