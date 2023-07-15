# Create and manage CI-CD-Target service.
resource "elestio_ci_cd_target" "demo_ci_cd_target" {
  project_id    = "2500"
  server_name   = "demo-ci_cd_target"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
