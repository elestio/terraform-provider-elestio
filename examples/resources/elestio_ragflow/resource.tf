resource "elestio_ragflow" "example" {
  project_id    = "2500"
  version       = "dev"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
