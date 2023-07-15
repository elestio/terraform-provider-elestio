# Create and manage Zulip service.
resource "elestio_zulip" "demo_zulip" {
  project_id    = "2500"
  server_name   = "demo-zulip"
  version       = "6.2-0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
