# Create and manage a Service.
resource "elestio_service" "myawesomeservice" {
  project_id    = "YOUR-PROJECT-ID"
  server_name   = "awesomeservice"
  server_type   = "SMALL-1C-2G"
  template_id   = 11 // postgreSQL
  version       = "14"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "YOUR-EMAIL"
}
