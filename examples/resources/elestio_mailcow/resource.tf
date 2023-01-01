# Create and manage MailCow service.
resource "elestio_mailcow" "my_mailcow" {
  project_id    = "2500"
  server_name   = "awesome-mailcow"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
