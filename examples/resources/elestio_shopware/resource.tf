resource "elestio_shopware" "example" {
  project_id    = "2500"
  version       = "latest-php8.2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
