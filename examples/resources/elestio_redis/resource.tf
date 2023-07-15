# Create and manage Redis service.
resource "elestio_redis" "demo_redis" {
  project_id    = "2500"
  server_name   = "demo-redis"
  version       = "6"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
