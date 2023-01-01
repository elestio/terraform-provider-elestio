# Create and manage Paperless-ngx service.
resource "elestio_paperless_ngx" "my_paperless_ngx" {
  project_id    = "2500"
  server_name   = "awesome-paperless_ngx"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
