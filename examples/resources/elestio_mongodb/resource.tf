resource "elestio_mongodb" "example" {
  project_id    = "2500"
  version       = "6"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
