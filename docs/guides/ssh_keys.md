---
page_title: "SSH Keys with Elestio Terraform Provider"
---

# SSH Keys with Elestio Terraform Provider

## What is an SSH Key?

SSH keys are a pair of cryptographic credentials (public and private key) used for secure, password-free authentication to remote servers. You submit your **public key** to Elestio, and use your **private key** to authenticate when connecting.

### Generate a Key

**Linux/macOS:**
```bash
ssh-keygen -t rsa
```

**Windows:**
1. Download [PuTTYgen](https://www.puttygen.com/)
2. Generate a new key (ed25519 or RSA recommended)
3. Save both private and public keys

## Quick Start

Use your public key in Terraform:

```tf
locals {
  key_data = provider::elestio::parse_ssh_key_data(file("~/.ssh/id_rsa.pub"))
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
```

Connect via SSH:
```bash
ssh -i ~/.ssh/id_rsa root@<service-hostname>
```

## Supported Key Types

Elestio supports the following SSH key types in OpenSSH format:
- `ssh-rsa`
- `ssh-ed25519`
- `ecdsa-sha2-nistp256`
- `ecdsa-sha2-nistp384`
- `ecdsa-sha2-nistp521`
- `ssh-dss`

SSH key comments must be stripped before use.

## Provider Functions

Elestio provides two functions to handle SSH key comments (requires Terraform 1.8+):

- **parse_ssh_key_data()** - Strips the comment from an SSH public key, returning only the key type and key data.
- **parse_ssh_key()** - Extracts both the key data and username from an SSH public key. Pass `null` to extract the username from the key comment.

### Example Usage

```tf
locals {
  # Option 1: Strip comment only, provide username manually
  key_data = provider::elestio::parse_ssh_key_data(file("~/.ssh/id_rsa.pub"))
  # Input:  "ssh-rsa AAAA... user@hostname"
  # Output: "ssh-rsa AAAA..."

  # Option 2: Extract both key and username from comment
  ssh_key = provider::elestio::parse_ssh_key(file("~/.ssh/id_rsa.pub"), null)

  # Option 3: Extract key but override username
  ssh_key_custom = provider::elestio::parse_ssh_key(file("~/.ssh/id_rsa.pub"), "custom-user")
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
    },
    local.ssh_key,        # Extracted username from comment
    local.ssh_key_custom, # Custom username
  ]
}
```