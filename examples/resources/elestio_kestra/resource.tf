# Create and manage Kestra service.
resource "elestio_kestra" "my_kestra" {
  project_id    = "2500"
  server_name   = "awesome-kestra"
  server_type   = "SMALL-1C-2G"
  version       = "develop-full"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
