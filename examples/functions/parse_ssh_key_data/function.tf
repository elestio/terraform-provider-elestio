# Strip the comment from an SSH public key
locals {
  key_data = provider::elestio::parse_ssh_key_data(file("~/.ssh/id_rsa.pub"))
  # Input:  "ssh-rsa AAAA... user@hostname"
  # Output: "ssh-rsa AAAA..."
}

resource "elestio_postgresql" "example" {
  project_id    = "project_id"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

  ssh_public_keys = [
    {
      username = "admin"
      key_data = local.key_data
    }
  ]
}
