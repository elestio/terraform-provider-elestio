# Create and manage ChangeDetection service.
resource "elestio_changedetection" "demo_changedetection" {
  project_id    = "2500"
  server_name   = "demo-changedetection"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
