---
page_title: 'Learn how to write your first Terraform configuration'
---

# Learn how to write your first Terraform configuration

In this guide, starting from an empty folder, you will learn how to deploy a Terraform infrastructure on Elestio !

If you have any doubt, you can check the final code of this guide in this repository [elestio-terraform-scratch][elestio-terraform-scratch](https://github.com/elestio-examples/elestio-terraform-scratch).

## Prepare the dependencies

-   [Sign up for Elestio if you haven't already](https://dash.elest.io/signup)
-   [Get your API token in the security settings page of your account ](https://dash.elest.io/account/security)
-   [Download and install Terraform](https://www.terraform.io/downloads)

You need a Terraform CLI version equal or higher than v0.14.0.
To ensure you're using the acceptable version of Terraform you may run the following command:

```bash
terraform -v
```

Your output should resemble:

```bash
Terraform v0.14.0 # any version >= v0.14.0 is OK
...
```

## Configure your project and services

Terraform files are used to define the structure and configuration of your infrastructure. It is generally a good idea to keep these definitions in **separate files** rather than combining them all in one file.

This section will explain how to organize a basic Terraform project :

1. **Create and move to an empty folder**

    -> **Note:** Here is an overview of the files we will create together. Do not create them yet.

    ```bash
    # do not create yet theses files
    |- main.tf
    |- variables.tf   # Defines the type of each variables required by main.tf
    |- secrets.tfvars # Contains the value of each variables defined in variables.tf
    ```

2. **Create a file `main.tf` and declare the provider adding the following lines :**

    ```hcl
    # main.tf

    terraform {
      required_providers {
        elestio = {
          source = "elestio/elestio"
          # check out the latest version available
          version = "0.10.0"
        }
      }
    }

    # Configure the Elestio Provider with your credentials
    provider "elestio" {
      email     = var.elestio_email
      api_token = var.elestio_api_token
    }
    ```

    As you can see, the email and API token are assigned to variables.  
    You should never put sensitive information directly in `.tf` files.

3. **Create a file `variables.tf` and declare variables adding the following lines :**

    ```hcl
    # variables.tf

    variable "elestio_email" {
        description = "Elestio Email"
        type        = string
    }

    variable "elestio_api_token" {
      description = "Elestio API Token"
      type        = string
      sensitive   = true
    }
    ```

    This file does not contain the values of these variables. We will have to declare them in another file.

4. Create a file `secrets.tfvars` and fill it with your values :

    ```bash
    # secrets.tfvars

    elestio_email      = "YOUR-EMAIL"
    elestio_api_token  = "YOUR-API-TOKEN"
    ```

    ~> **Note:** Do not commit with Git this file ! Sensitive information such as an API token should never be pushed. For more information on how to securely authenticate, please read the [authentication documentation](https://registry.terraform.io/providers/elestio/elestio/latest/docs#authentication).

5. **Go back to the `main.tf` file and add the following lines :**

    ```hcl
    # add to main.tf

    # Create a Project
    resource "elestio_project" "pg_project" {
      name             = "PostgreSQL Project"
      description      = "Contains a postgres database"
      technical_emails = var.elestio_email
    }
    ```

    To contain our PostgreSQL service, we will have to create a new project on Elestio.  
    Instead of using the web interface, we can also declare it via terraform.

6. **In the same `main.tf` file, add the following lines :**

    ```terraform
    # add to main.tf

    # Create a PostgreSQL Service
    resource "elestio_postgresql" "my_service" {
      project_id       = elestio_project.my_project.id
      provider_name    = "hetzner"
      datacenter       = "fsn1"
      server_type      = "SMALL-1C-2G"
      firewall_enabled = true
    }

    # Output the command that can be used to connect to the database
    output "pg_service_psql_command" {
      value       = elestio_postgresql.my_service.database_admin.command
      description = "The PSQL command to connect to the database."
      sensitive   = true
    }
    ```

    Terraform takes care of managing the dependencies and creating the different resources in the right order. As you can see, `project_id` will be filled with the value of the Project Resource that will be created.

## Apply the Terraform configuration

1. **Download and install the Elestio provider defined in the configuration :**

    ```bash
    terraform init
    ```

2. **Ensure the configuration is syntactically valid and internally consistent:**

    ```bash
    terraform validate
    ```

3. **Apply the configuration :**

    ```bash
    terraform apply -var-file="secrets.tfvars"
    ```

    Deployment time varies by service, provider, datacenter and server type.

4. **Voila, you have created a Project and PostgreSQL Service using Terraform !**

    You can visit the [Elestio web dashboard](https://dash.elest.io/) to see these ressources.

## Update the configuration

1. Change the `firewall_enabled` value to false in `main.tf` and run the following command :

    ```bash
    terraform apply -var-file="secrets.tfvars"
    ```

    This will update the configuration and destroy the firewall.

2. Revert the change in `main.tf` and run the following command :

    ```bash
    terraform apply -var-file="secrets.tfvars"
    ```

    This will update the configuration and create the firewall again.

Some changes (ex: `datacenter`) require the creation of new resources and the destruction of old resources.
Terraform will show you the resources to be created and destroyed before prompting you to confirm.
You can loose data if you destroy a resource, so be careful.

## (Optional) Access to the database

Let's try to connect to the database to see if everything worked well

First, you need to [install psql. ](https://www.timescale.com/blog/how-to-install-psql-on-mac-ubuntu-debian-windows/)

After that, run this command :

```bash
eval "$(terraform output -raw pg_service_psql_command)"
```

-> **Note:** The command to leave psql terminal is `\q`

## Clean up

Run the following command to destroy all the resources you created:

```bash
terraform destroy -var-file="secrets.tfvars"
```

This command destroys all the resources specified in your Terraform state. `terraform destroy` doesn't destroy resources running elsewhere that aren't managed by the current Terraform project.

Now you've created and destroyed an entire Elestio deployment!

Visit the [Elestio Dashboard](https://dash.elest.io/) to verify the resources have been destroyed to avoid unexpected charges.
