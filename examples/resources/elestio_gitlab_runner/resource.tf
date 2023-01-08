# Create and manage Gitlab-runner service.
resource "elestio_gitlab_runner" "my_gitlab_runner" {
  project_id    = "2500"
  server_name   = "awesome-gitlab_runner"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
