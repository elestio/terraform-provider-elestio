# Create and manage Trilium service.
resource "elestio_trilium" "my_trilium" {
  project_id    = "2500"
  server_name   = "awesome-trilium"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
