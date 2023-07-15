# Create and manage LimeSurvey service.
resource "elestio_limesurvey" "demo_limesurvey" {
  project_id    = "2500"
  server_name   = "demo-limesurvey"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
