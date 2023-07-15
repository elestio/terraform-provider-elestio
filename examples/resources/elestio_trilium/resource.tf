# Create and manage Trilium service.
resource "elestio_trilium" "demo_trilium" {
  project_id    = "2500"
  server_name   = "demo-trilium"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
