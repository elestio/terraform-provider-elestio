# Create and manage Formbricks service.
resource "elestio_formbricks" "demo_formbricks" {
  project_id    = "2500"
  server_name   = "demo-formbricks"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
