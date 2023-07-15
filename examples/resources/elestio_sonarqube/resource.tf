# Create and manage SonarQube service.
resource "elestio_sonarqube" "demo_sonarqube" {
  project_id    = "2500"
  server_name   = "demo-sonarqube"
  version       = "9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
