# Create and manage Penpot service.
resource "elestio_penpot" "demo_penpot" {
  project_id    = "2500"
  server_name   = "demo-penpot"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
