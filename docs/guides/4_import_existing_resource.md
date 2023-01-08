---
page_title: "Import an existing resource"
---

# Import an existing resource

You can use the `terraform import` command to import in the Elestio state an existing project or service already running.

## Project

```shell
# Import a project by specifying the project ID.
terraform import elestio_project.myawesomeproject project_id
```

1. Declare the resource in your terraform file

   ```tf
   resource "elestio_project" "example_project" {
    name = "example-project"
   # ...
   }
   ```

2. Retrieve your ProjectID on the [projects list page](https://dash.elest.io/projects).

   ![projects list page](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-02-01-12.png)

3. Execute the import command

   ```shell
   terraform import elestio_project.example_project 2434
   ```

   Then you can run a `terraform apply` command and elestio will handle your project resource as an update and not a new resource to create.

## Service

```shell
# Import a service by specifying the Project ID it belongs to, and the service ID (spaced by a comma).
terraform import elestio_service.myawesomeservice project_id,service_id
```

1. Declare the resource in your terraform file

   ```tf
   resource "elestio_postgres" "example_postgres" {
    project_id    = elestio_project.example_project.id
    # ...
   }
   ```

2. Retrieve your ProjectID on the [projects list page](https://dash.elest.io/projects).

3. Retrieve your ServiceID on the service page.

   ![service page](https://docs.elest.io/uploads/images/gallery/2023-01/scaled-1680-/screenshot-2023-01-08-at-02-04-26.png)

4. Execute the import command

   ```shell
   terraform import elestio_postgres.example_postgres 2434,27265610
   ```

   Then you can run a `terraform apply` command and elestio will handle your service resource as an update and not a new resource to create.
