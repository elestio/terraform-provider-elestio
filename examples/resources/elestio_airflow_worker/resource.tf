# Create and manage Airflow-worker service.
resource "elestio_airflow_worker" "demo_airflow_worker" {
  project_id    = "2500"
  server_name   = "demo-airflow_worker"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
