---
page_title: "Providers, datacenters and server types"
---

# Providers, datacenters and server types

This guide will list the available options for `provider_name`, `datacenter` and `server_type` variables when you want to manage a service resource with terraform.

```tf
resource "elestio_vault" "my_vault" {
  // ..
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "SMALL-1C-2G"
  // ...
}
```

As this information can be updated often, we cannot put a fixed list in this documentation.
You will learn how to get this information from the elestio website.

## Listing all options

When you create a service via the website, all three pieces of information (providers, datacenters and server types) are listed on a single page.

1. Sign In on the [Elestio Web Dashboard](https://dash.elest.io/).

2. Click on the button **Create a new service** (or **Deploy my first service**).

![click buttonc create](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-01-15-44.png)

3. Select one of the services (the choice does not matter)

![select a service](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-01-16-20.png)

4. You arrive on this page :
   ![providers list](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/1beimage.png)

   All providers are listed.

   When you select one, you will find the list of datacenters below.

   Finally after selecting a datacenter, you will find the list of server types below.

## Providers

![providers](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-01-20-12.png)

For each of the providers, here are the values to pass to terraform **provider_name**:

- Hetzner Cloud -> `hetzner`
- Digital Ocean -> `do`
- Amazon Lightsail -> `lightsail`
- Linode Cloud -> `linode`
- Vultr Cloud -> `vultr`
- Scaleway -> `scaleway`

## Datacenters

![datacenters](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-01-19-31.png)

You can pass **datacenter** variable in as it is written on the site. For example: `fsn1` (Hetzner) or `eu-central-1` (Amazon Lightsail).

## Server types

![server types](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-01-24-56.png)

You can pass **server_type** variable in as it is written on the site. For example: `SMALL-1C-2G` (Hetzner) or `MICRO-1C-1G` (Amazon Lightsail).

-> **Info:** Note that after creating a service, you can only modify a server type with a more powerful one. Downgrade is not allowed.

## Example

If I want to host my service on **Amazon Lightsail** in **Ireland** (Dublin) with a **1 CPU, 1 GB RAM, 2 TB Bandwith** plan :

```tf
resource "elestio_vault" "my_vault" {
  // ..
  provider_name = "lightsail"
  datacenter    = "eu-west-1"
  server_type   = "MICRO-1C-1G"
  // ...
}
```
