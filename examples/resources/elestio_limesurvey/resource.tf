# Create and manage LimeSurvey service.
resource "elestio_limesurvey" "my_limesurvey" {
  project_id    = "2500"
  server_name   = "awesome-limesurvey"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
