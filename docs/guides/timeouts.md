---
page_title: "Operation Timeouts with Elestio Terraform Provider"
---

# Operation Timeouts

When Terraform creates a service, it waits for the service to fully start before
continuing. Some services take longer than others depending on the template,
server type, and datacenter.

By default the provider waits **20 minutes** for each create, update, and delete
operation. If that isn't enough, you'll see an error like this even though the
service is actually running on Elestio:

```
Error Creating Service: Timed out after 20m0s while waiting for the service to start
```

## Setting a custom timeout

Add a `timeouts` block to the resource and set a longer duration (e.g. `"30m"`,
`"1h"`):

```tf
resource "elestio_supabase" "example" {
  project_id    = "project_id"
  provider_name = "linode"
  datacenter    = "ca-central"
  server_type   = "LARGE-4C-8G"

  timeouts {
    create = "30m"
  }
}
```

You can set `create`, `update`, and `delete` independently. Any value you leave
out keeps the 20-minute default. This block works on every `elestio_*` service
resource.
