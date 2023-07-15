# Create and manage Gitlab service.
resource "elestio_gitlab" "demo_gitlab" {
  project_id    = "2500"
  server_name   = "demo-gitlab"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
