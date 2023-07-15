# Create and manage RabbitMQ service.
resource "elestio_rabbitmq" "demo_rabbitmq" {
  project_id    = "2500"
  server_name   = "demo-rabbitmq"
  version       = "3-management"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
