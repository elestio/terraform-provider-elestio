resource "elestio_shopware" "example" {
  project_id    = "2500"
  version       = "6.4.20-php8.2"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
