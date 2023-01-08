# Create and manage Kafka service.
resource "elestio_kafka" "my_kafka" {
  project_id    = "2500"
  server_name   = "awesome-kafka"
  server_type   = "SMALL-1C-2G"
  version       = "5.5.5"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
