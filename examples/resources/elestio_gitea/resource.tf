# Create and manage Gitea service.
resource "elestio_gitea" "demo_gitea" {
  project_id    = "2500"
  server_name   = "demo-gitea"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
