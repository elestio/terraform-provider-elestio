resource "elestio_microk8s" "example" {
  project_id    = "2500"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
