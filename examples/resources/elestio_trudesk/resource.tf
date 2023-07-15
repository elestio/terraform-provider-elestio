# Create and manage Trudesk service.
resource "elestio_trudesk" "demo_trudesk" {
  project_id    = "2500"
  server_name   = "demo-trudesk"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
