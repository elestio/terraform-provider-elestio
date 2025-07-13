resource "elestio_minio" "example" {
  project_id    = "2500"
  version       = "RELEASE.2025-04-22T22-12-26Z"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
