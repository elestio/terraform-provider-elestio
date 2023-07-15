# Create and manage Paperless-ngx service.
resource "elestio_paperless_ngx" "demo_paperless_ngx" {
  project_id    = "2500"
  server_name   = "demo-paperless_ngx"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
