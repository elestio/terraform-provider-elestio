# Create and manage Typesense service.
resource "elestio_typesense" "demo_typesense" {
  project_id    = "2500"
  server_name   = "demo-typesense"
  version       = "0.23.0.rc66"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
