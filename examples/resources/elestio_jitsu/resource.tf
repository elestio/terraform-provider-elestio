# Create and manage Jitsu service.
resource "elestio_jitsu" "my_jitsu" {
  project_id    = "2500"
  server_name   = "awesome-jitsu"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
