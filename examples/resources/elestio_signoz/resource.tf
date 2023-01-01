# Create and manage SigNoz service.
resource "elestio_signoz" "my_signoz" {
  project_id    = "2500"
  server_name   = "awesome-signoz"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
