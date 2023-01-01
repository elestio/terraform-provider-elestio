# Create and manage Redis service.
resource "elestio_redis" "my_redis" {
  project_id    = "2500"
  server_name   = "awesome-redis"
  server_type   = "SMALL-1C-2G"
  version       = "6"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
