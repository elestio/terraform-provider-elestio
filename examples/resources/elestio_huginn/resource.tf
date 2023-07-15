# Create and manage Huginn service.
resource "elestio_huginn" "demo_huginn" {
  project_id    = "2500"
  server_name   = "demo-huginn"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
