# Create and manage SonarQube service.
resource "elestio_sonarqube" "my_sonarqube" {
  project_id    = "2500"
  server_name   = "awesome-sonarqube"
  server_type   = "SMALL-1C-2G"
  version       = "9"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "exemple@mail.com"
}
