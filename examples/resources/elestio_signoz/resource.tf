# Create and manage SigNoz service.
resource "elestio_signoz" "demo_signoz" {
  project_id    = "2500"
  server_name   = "demo-signoz"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
