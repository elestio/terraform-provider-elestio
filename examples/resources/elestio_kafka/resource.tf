# Create and manage Kafka service.
resource "elestio_kafka" "demo_kafka" {
  project_id    = "2500"
  server_name   = "demo-kafka"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
