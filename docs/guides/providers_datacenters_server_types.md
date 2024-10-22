---
page_title: "Providers, datacenters and server types"
---

# Providers, datacenters and server types

This guide explain how to find available options for `provider_name`, `datacenter` and `server_type` variables when you want to manage a service resource with terraform.

```tf
resource "elestio_vault" "my_vault" {
  // ..
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
  // ...
}
```

As this information can be updated often, we cannot put a fixed list in this documentation.
You will learn how to get this information from the elestio website.

## Listing all options

When you create a service via the website, all three pieces of information (**providers**, **data centers**, and **server types**) **are listed on a single page**. You can copy the configuration from there and paste it into your Terraform file.

[![1- Navigate to Elestio website](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-55-15.png)](https://elest.io/)

![2- Login to the dashboard](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-55-32.png)

![3- Click on the button Deploy my first service](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-55-53.png)

![4- Search service by name](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-56-15.png)

![5- Select the service](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-56-27.png)

![6- Choose the provider](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-56-51.png)

![7- Choose a datacenter](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-57-04.png)

![8- Choose a server type](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-57-20.png)

![9- Select a software version](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-57-39.png)

![10- Click Copy Terraform Config](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-57-50.png)

![11- Copy the config and paste it in your terraform file](https://docs.elest.io/uploads/images/gallery/2023-10/scaled-1680-/cleanshot-2023-10-03-at-12-58-05.png)
