resource "elestio_rke2" "example" {
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
      "port"     = "80"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "6443"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "9345"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "2379"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "10250"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "2380"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "8472"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "30000-32767"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "10259"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "10257"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "10256"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    }
  ]
}
