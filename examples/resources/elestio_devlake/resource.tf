resource "elestio_devlake" "example" {
  project_id    = "2500"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
