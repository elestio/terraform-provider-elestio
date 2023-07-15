# Create and manage LocalStack service.
resource "elestio_localstack" "demo_localstack" {
  project_id    = "2500"
  server_name   = "demo-localstack"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
