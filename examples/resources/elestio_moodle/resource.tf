resource "elestio_moodle" "example" {
  project_id    = "2500"
  version       = "4.2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
