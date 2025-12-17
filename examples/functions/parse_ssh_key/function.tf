# Parse SSH keys
locals {
  # Extract username from SSH key comment (pass null explicitly)
  ssh_key = provider::elestio::parse_ssh_key(file("~/.ssh/id_rsa.pub"), null)

  # Override username with a custom value
  ssh_key_custom = provider::elestio::parse_ssh_key(file("~/.ssh/id_rsa.pub"), "custom-username")
}

# Use the parsed SSH keys in a resource
resource "elestio_postgresql" "example" {
  project_id    = "project_id"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

  ssh_public_keys = [
    local.ssh_key,        # { key_data = "ssh-rsa AAAA...", username = "user@hostname" }
    local.ssh_key_custom, # { key_data = "ssh-rsa AAAA...", username = "custom-username" }
  ]
}
