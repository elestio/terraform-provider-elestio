# Create and manage Airflow-worker service.
resource "elestio_airflow_worker" "my_airflow_worker" {
  project_id    = "2500"
  server_name   = "awesome-airflow_worker"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
