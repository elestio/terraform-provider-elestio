---
page_title: "Firewall Configuration"
description: |-
  Learn how to configure firewall rules for Elestio service resources.
---

# Firewall Configuration

This guide explains how to configure firewall rules for Elestio service resources using Terraform.

See also the [dashboard documentation](https://docs.elest.io/books/security/page/network-firewall) for more information.

## Overview

Firewall is enabled by default on all service resources and controls network access. You can define custom rules to allow or block access to specific ports, or disable the firewall entirely if needed.

## Basic Configuration

Rules are defined in the `firewall_user_rules` attribute:

```terraform
resource "elestio_postgresql" "database" {
  project_id    = "2500"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

  # This is the default value if you don't specify it.
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
    # Application-specific ports
    {
      "type"     = "input"
      "port"     = "80"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "443"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "25432"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    }
  ]
}
```

### Rule Attributes

- `port` (Required) - Single port (`"22"`) or port range (`"400-410"`)
- `protocol` (Required) - Either `tcp` or `udp`
- `type` (Required) - Rule type: `input` (currently the only supported value)
- `targets` (Required) - CIDR block(s) for allowed sources. Use `["0.0.0.0/0", "::/0"]` to allow access from all IPv4 and IPv6 addresses

## Default Rules

The `firewall_user_rules` attribute is optional. When omitted, it includes:

**Required system ports** (automatically included):
- `22/tcp` - SSH access for Elestio automation (must have targets `["0.0.0.0/0", "::/0"]`)
- `4242/udp` - Nebula VPN for Global IP connectivity (must have targets `["0.0.0.0/0", "::/0"]`)
- `80/tcp` - Let's Encrypt (when using custom domains, must have targets `["0.0.0.0/0", "::/0"]`)

!> **Important** System ports **22/tcp**, **4242/udp**, and **80/tcp** (when using custom domains) must have **exactly** targets `["0.0.0.0/0", "::/0"]` - no additional targets allowed. These ports are critical for Elestio automation, Global IP connectivity, and SSL certificate generation.

**Application-specific ports** (varies by service):
- PostgreSQL: `80/tcp`, `443/tcp`, `25432/tcp` (note: not the standard 5432)
- Redis: `26379/tcp`, `80/tcp`, `443/tcp`, `26380/tcp` (note: not the standard 6379)
- Web apps: `80/tcp`, `443/tcp`

~> **Warning** Elestio may use different ports than official software defaults. For example, PostgreSQL uses port **25432** instead of **5432**. Check each resource's documentation for its specific ports.

To view actual default rules for your service:
- Check the resource documentation for the `firewall_user_rules` attribute description
- Comment out `firewall_user_rules` and run `terraform plan`

## Disabling the Firewall

!> **Important** When disabling the firewall, you **must** explicitly set `firewall_user_rules = []`:

```terraform
resource "elestio_postgresql" "database" {
  project_id    = elestio_project.my_project.id
  server_name   = "my-database"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"

  firewall_enabled    = false  # Not recommended
  firewall_user_rules = []     # Required
}
```

## Dashboard Tool Ports

Dashboard tools automatically add their ports to the firewall:

- **VS Code** - `18345/tcp`
- **Open Terminal** - `18374/tcp`
- **File Explorer** - `18346/tcp`
- **Tail Logs** - `18445/tcp`
- **Terminal** - `18344/tcp`

Terraform ignores these in `firewall_user_rules`. They appear in the read-only `firewall_tool_rules` attribute.

-> **Tip** To keep a tool port always open, add it explicitly to `firewall_user_rules`.

### Removing Tool Ports

Close unused tool ports with `firewall_remove_tool_ports`:

```terraform
resource "elestio_postgresql" "database" {
  # ... other configuration ...
  
  firewall_remove_tool_ports = true
}
```

!> **Important** This is a one-time cleanup action. After running `terraform apply`, set it back to `false`. The cleanup only executes when transitioning from `false` to `true`. Tools will still add their ports when used again.

## Viewing Rules

View tool ports:
```terraform
output "tool_ports" {
  value = elestio_postgresql.database.firewall_tool_rules
}
```

## Example: Restricting Database Access

Limit PostgreSQL access to specific backend servers:

```terraform
resource "elestio_postgresql" "database" {
  project_id    = elestio_project.my_project.id
  server_name   = "secure-db"
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"
  
  firewall_user_rules = [
    # Required system ports (must be exactly these targets)
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
    # Application-specific ports
    {
      "type"     = "input"
      "port"     = "80"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "443"
      "protocol" = "tcp"
      "targets"  = ["0.0.0.0/0", "::/0"]
    },
    {
      "type"     = "input"
      "port"     = "25432"
      "protocol" = "tcp"
      "targets"  = ["212.11.12.13/32", "212.11.12.14/32"] # Restricted to specific backend servers
    }
  ]
}
```