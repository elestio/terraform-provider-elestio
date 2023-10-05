# Create and manage TotaljsFlow service.
resource "elestio_totaljsflow" "demo_totaljsflow" {
  project_id    = "2500"
  server_name   = "demo-totaljsflow"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
