resource "elestio_zulip" "example" {
  project_id    = "2500"
  version       = "6.2-0"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
