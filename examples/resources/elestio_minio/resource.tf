# Create and manage MinIO service.
resource "elestio_minio" "demo_minio" {
  project_id    = "2500"
  server_name   = "demo-minio"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
