# Create and manage Gitea service.
resource "elestio_gitea" "my_gitea" {
  project_id    = "2500"
  server_name   = "awesome-gitea"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
