# Create and manage Docspell service.
resource "elestio_docspell" "demo_docspell" {
  project_id    = "2500"
  server_name   = "demo-docspell"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
