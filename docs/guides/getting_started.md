---
page_title: 'Getting started with the Elestio provider'
---

# Getting Started with the Elestio Provider

In this tutorial, we will discover together the Elestio Terraform Provider.

As a goal, we will deploy not one but two PostgreSQL services!
And to complicate things, they will be hosted in different regions.

You don't have to code anything, you can clone this repository [elestio-terraform-postgres-regions](https://github.com/elestio-examples/elestio-terraform-postgres-regions) that contain required code for this guide.

## Signup

Start by registering on the Elestio website: https://dash.elest.io/signup

-> You will need to load some credits: https://dash.elest.io/account/payment
<br>Don't worry, Elestio services are charged by the hour so this example won't cost much.

## Get API token

You will need your API Token to identify yourself to the provider.
You can find this token in the security settings of your account: https://dash.elest.io/account/security

## Install Terraform CLI

You need a version equal or higher than v0.14.0.

You can find some documentation [here](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli#install-terraform).

To ensure you're using the acceptable version of Terraform you may run the following command:

```bash
$ terraform -v
```

Your output should resemble:

```bash
Terraform v0.14.0 # any version >= v0.14.0 is OK
...
```

## Deploy

1.  Clone the [repository](https://github.com/elestio-examples/elestio-terraform-postgres-regions) containing the example configurations:

    ```bash
    git clone https://github.com/elestio-examples/elestio-terraform-postgres-regions.git
    ```

    ```bash
    |- main.tf
    |- variables.tf   # Defines the type of each variables required by main.tf
    |- secrets.tfvars # Contains the value of each variables defined in variables.tf
    ```

2.  Rename `secrets.tfvars.tmp` file `secrets.tfvars` and fill in the appropriate values:

    ```bash
    mv secrets.tfvars.tmp secrets.tfvars
    ```

    ```terraform
    # file secrets.tfvars

    elestio_email     = "<elestio_email>"
    elestio_api_token = "<elestio_api_token>"
    ```

3.  Download and install the Elestio provider defined in the configuration:

    ```bash
    terraform init
    ```

4.  Ensure the configuration is syntactically valid and internally consistent:

    ```bash
    terraform validate
    ```

5.  Apply the configuration:

    ```bash
    terraform apply -var-file="secrets.tfvars"
    ```

    Deployment time varies by service, provider, datacenter and server type.

6.  You have created one Project and two PostgreSQL Services using Terraform! Visit the [Elestio Dashboard](https://dash.elest.io/) to see these resources.

## (Optional) Access databases

You need to [install psql](https://www.timescale.com/blog/how-to-install-psql-on-mac-ubuntu-debian-windows/).

Run these two commands separately:

```bash
eval "$(terraform output -raw pg_europe_psql)"
...
eval "$(terraform output -raw pg_asia_psql)"
...
```

-> **Note:** The command to leave psql terminal is `\q`.

## Clean up

Run the following command to destroy all the resources you created:

```bash
terraform destroy -var-file="secrets.tfvars"
```

This command destroys all the resources specified in your Terraform state. `terraform destroy` doesn't destroy resources running elsewhere that aren't managed by the current Terraform project.

Now you've created and destroyed an entire Elestio deployment!

Visit the [Elestio Dashboard](https://dash.elest.io/) to verify the resources have been destroyed to avoid unexpected charges.
