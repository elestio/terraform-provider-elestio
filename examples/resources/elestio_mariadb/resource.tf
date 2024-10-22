resource "elestio_mariadb" "example" {
  project_id    = "2500"
  version       = "10.9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
