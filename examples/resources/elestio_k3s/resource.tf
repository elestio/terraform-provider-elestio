# Create and manage K3S service.
resource "elestio_k3s" "demo_k3s" {
  project_id    = "2500"
  server_name   = "demo-k3s"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
