# Create and manage SFTPGo service.
resource "elestio_sftpgo" "my_sftpgo" {
  project_id    = "2500"
  server_name   = "awesome-sftpgo"
  server_type   = "SMALL-1C-2G"
  version       = "latest"
  provider_name = "hetzner"
  datacenter    = "fsn1"
  support_level = "level1"
  admin_email   = "example@mail.com"
  ssh_keys      = []
}
