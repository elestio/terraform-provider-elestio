resource "elestio_ubuntu" "example" {
  project_id    = "2500"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
