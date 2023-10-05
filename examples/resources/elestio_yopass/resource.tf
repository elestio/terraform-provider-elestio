# Create and manage Yopass service.
resource "elestio_yopass" "demo_yopass" {
  project_id    = "2500"
  server_name   = "demo-yopass"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
