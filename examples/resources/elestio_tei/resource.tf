resource "elestio_tei" "example" {
  project_id    = "2500"
  version       = "cpu-latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
