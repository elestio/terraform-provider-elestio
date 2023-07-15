# Create and manage MailCow service.
resource "elestio_mailcow" "demo_mailcow" {
  project_id    = "2500"
  server_name   = "demo-mailcow"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
