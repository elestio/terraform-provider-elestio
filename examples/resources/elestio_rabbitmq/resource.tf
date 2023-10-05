resource "elestio_rabbitmq" "example" {
  project_id    = "2500"
  version       = "3-management"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
