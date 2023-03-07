# Create and manage RabbitMQ service.
resource "elestio_rabbitmq" "my_rabbitmq" {
  project_id    = "2500"
  server_name   = "awesome-rabbitmq"
  server_type   = "SMALL-1C-2G"
  version       = "3-management"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
