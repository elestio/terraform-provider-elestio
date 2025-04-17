resource "elestio_mautic" "example" {
  project_id    = "2500"
  version       = "5.2.3-fpm"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
