# Create and manage Zabbix service.
resource "elestio_zabbix" "demo_zabbix" {
  project_id    = "2500"
  server_name   = "demo-zabbix"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
