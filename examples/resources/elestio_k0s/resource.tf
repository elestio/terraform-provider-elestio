# Create and manage k0s service.
resource "elestio_k0s" "demo_k0s" {
  project_id    = "2500"
  server_name   = "demo-k0s"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
