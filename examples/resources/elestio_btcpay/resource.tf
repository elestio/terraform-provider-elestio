# Create and manage BTCPay service.
resource "elestio_btcpay" "demo_btcpay" {
  project_id    = "2500"
  server_name   = "demo-btcpay"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
