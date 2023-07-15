# Create and manage a Service.
resource "elestio_service" "demo_service" {
  project_id    = "YOUR-PROJECT-ID"
  template_id   = 11 // postgreSQL
  server_name   = "awesomeservice"
  version       = "14"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
