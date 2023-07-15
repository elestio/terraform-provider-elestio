# Create and manage Airflow service.
resource "elestio_airflow" "demo_airflow" {
  project_id    = "2500"
  server_name   = "demo-airflow"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
