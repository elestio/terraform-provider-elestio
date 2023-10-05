resource "elestio_itop" "example" {
  project_id    = "2500"
  version       = "latest-base"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
