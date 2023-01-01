# Create and manage CI-CD-Target service.
resource "elestio_ci_cd_target" "my_ci_cd_target" {
  project_id    = "2500"
  server_name   = "awesome-ci_cd_target"
  server_type   = "SMALL-1C-2G"
  version       = ""
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
