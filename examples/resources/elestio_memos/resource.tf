# Create and manage Memos service.
resource "elestio_memos" "demo_memos" {
  project_id    = "2500"
  server_name   = "demo-memos"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
