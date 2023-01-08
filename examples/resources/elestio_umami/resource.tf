# Create and manage Umami service.
resource "elestio_umami" "my_umami" {
  project_id    = "2500"
  server_name   = "awesome-umami"
  server_type   = "SMALL-1C-2G"
  version       = "postgresql-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
