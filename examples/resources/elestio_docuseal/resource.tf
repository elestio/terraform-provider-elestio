# Create and manage DocuSeal service.
resource "elestio_docuseal" "demo_docuseal" {
  project_id    = "2500"
  server_name   = "demo-docuseal"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
