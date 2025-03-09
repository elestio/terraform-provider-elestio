resource "elestio_revolt" "example" {
  project_id    = "2500"
  version       = "20241024-1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
