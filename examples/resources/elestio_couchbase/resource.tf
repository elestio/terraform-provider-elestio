# Create and manage Couchbase service.
resource "elestio_couchbase" "my_couchbase" {
  project_id    = "2500"
  server_name   = "awesome-couchbase"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
