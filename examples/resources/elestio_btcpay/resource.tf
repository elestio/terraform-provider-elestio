# Create and manage BTCPay service.
resource "elestio_btcpay" "my_btcpay" {
  project_id    = "2500"
  server_name   = "awesome-btcpay"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
