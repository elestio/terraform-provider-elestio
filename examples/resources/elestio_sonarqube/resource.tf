resource "elestio_sonarqube" "example" {
  project_id    = "2500"
  version       = "9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
