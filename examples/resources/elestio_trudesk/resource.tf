# Create and manage Trudesk service.
resource "elestio_trudesk" "my_trudesk" {
  project_id    = "2500"
  server_name   = "awesome-trudesk"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
