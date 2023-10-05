# Create and manage Cal.com service.
resource "elestio_cal_com" "demo_cal_com" {
  project_id    = "2500"
  server_name   = "demo-cal_com"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
