resource "elestio_easyappointments" "example" {
  project_id    = "2500"
  version       = "1.4.3"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
