resource "elestio_columnstore" "example" {
  project_id    = "2500"
  version       = "latest"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

  # Default firewall rules
  firewall_user_rules = [
    # Required system ports
    {
      "type"     = "input"
      "port"     = "22"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "4242"
      "protocol" = "udp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    # Application ports
    {
      "type"     = "input"
      "port"     = "23060"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "443"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    }
  ]
}
