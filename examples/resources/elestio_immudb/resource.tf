# Create and manage ImmuDB service.
resource "elestio_immudb" "demo_immudb" {
  project_id    = "2500"
  server_name   = "demo-immudb"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
