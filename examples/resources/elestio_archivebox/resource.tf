# Create and manage ArchiveBox service.
resource "elestio_archivebox" "demo_archivebox" {
  project_id    = "2500"
  server_name   = "demo-archivebox"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
