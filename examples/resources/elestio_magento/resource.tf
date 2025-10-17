resource "elestio_magento" "example" {
  project_id    = "2500"
  version       = "8.4-fpm-1.4.4"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
