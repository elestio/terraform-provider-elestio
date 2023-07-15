# Create and manage Microk8s service.
resource "elestio_microk8s" "demo_microk8s" {
  project_id    = "2500"
  server_name   = "demo-microk8s"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
