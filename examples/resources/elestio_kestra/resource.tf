# Create and manage Kestra service.
resource "elestio_kestra" "demo_kestra" {
  project_id    = "2500"
  server_name   = "demo-kestra"
  version       = "develop-full"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
