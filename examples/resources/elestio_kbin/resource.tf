# Create and manage KBIN service.
resource "elestio_kbin" "demo_kbin" {
  project_id    = "2500"
  server_name   = "demo-kbin"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
