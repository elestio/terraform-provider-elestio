# Create and manage FlowiseAI service.
resource "elestio_flowiseai" "demo_flowiseai" {
  project_id    = "2500"
  server_name   = "demo-flowiseai"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
