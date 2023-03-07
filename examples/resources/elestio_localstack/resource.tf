# Create and manage LocalStack service.
resource "elestio_localstack" "my_localstack" {
  project_id    = "2500"
  server_name   = "awesome-localstack"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
