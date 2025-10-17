resource "elestio_revolt" "example" {
  project_id    = "2500"
  version       = "20250210-1"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
