# Create and manage iTop service.
resource "elestio_itop" "demo_itop" {
  project_id    = "2500"
  server_name   = "demo-itop"
  version       = "latest-base"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
