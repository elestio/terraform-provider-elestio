# Create and manage Airflow service.
resource "elestio_airflow" "my_airflow" {
  project_id    = "2500"
  server_name   = "awesome-airflow"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
}
