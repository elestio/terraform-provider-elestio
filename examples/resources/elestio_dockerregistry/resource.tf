# Create and manage DockerRegistry service.
resource "elestio_dockerregistry" "my_dockerregistry" {
  project_id    = "2500"
  server_name   = "awesome-dockerregistry"
  server_type   = "SMALL-1C-2G"
  version       = "2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
