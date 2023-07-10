# Create and manage PeerTube service.
resource "elestio_peertube" "my_peertube" {
  project_id    = "2500"
  server_name   = "awesome-peertube"
  server_type   = "SMALL-1C-2G"
  version       = "production-bullseye"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
