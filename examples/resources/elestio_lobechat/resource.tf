resource "elestio_lobechat" "example" {
  project_id    = "2500"
  version       = "0.1.16"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
