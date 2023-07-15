# Create and manage SFTPGo service.
resource "elestio_sftpgo" "demo_sftpgo" {
  project_id    = "2500"
  server_name   = "demo-sftpgo"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
}
