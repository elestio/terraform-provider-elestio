# Create and manage ApacheNiFi service.
resource "elestio_apachenifi" "demo_apachenifi" {
  project_id    = "2500"
  server_name   = "demo-apachenifi"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
