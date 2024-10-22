resource "elestio_rabbitmq" "example" {
  project_id    = "2500"
  version       = "3-management"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
