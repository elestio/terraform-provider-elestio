# Create and manage Nomad service.
resource "elestio_nomad" "my_nomad" {
  project_id    = "2500"
  server_name   = "awesome-nomad"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
