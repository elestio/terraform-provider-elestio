# Create and manage GrowthBook service.
resource "elestio_growthbook" "demo_growthbook" {
  project_id    = "2500"
  server_name   = "demo-growthbook"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
