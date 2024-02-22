resource "elestio_pritunl" "example" {
  project_id    = "2500"
  version       = "1.32.3746.95"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
