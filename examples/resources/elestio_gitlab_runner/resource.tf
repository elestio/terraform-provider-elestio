# Create and manage Gitlab-runner service.
resource "elestio_gitlab_runner" "demo_gitlab_runner" {
  project_id    = "2500"
  server_name   = "demo-gitlab_runner"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
