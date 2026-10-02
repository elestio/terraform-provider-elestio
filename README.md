# Elestio Terraform Provider

The Terraform provider for [Elest.io](https://elest.io/).
Elestio is a fully managed DevOps platform to deploy your code and open-source software.

**See the [official documentation](https://registry.terraform.io/providers/elestio/elestio/latest/docs) to learn about all the possible services and resources.**

## Get Started

Let's deploy a **PostgreSQL** database in a few minutes.

- [Signup for Elestio](https://dash.elest.io/signup)
- [Get your API Token](https://dash.elest.io/account/security)
- Create a file named `main.tf` with the content below:

```hcl
terraform {
  required_providers {
    elestio = {
      source  = "elestio/elestio"
    }
  }
}

# Authenticate
provider "elestio" {
  email = "your-account-email"
  api_token = "your-api-token"
}

# Project that will contain the postgres service
resource "elestio_project" "project" {
  name             = "Demo"
}

# Service postgres
resource "elestio_postgresql" "postgres" {
  project_id    = elestio_project.project.id
  provider_name = "netcup"
  datacenter    = "nbg"
  server_type   = "MEDIUM-2C-4G"
}

# Retrieve the command to access the database
output "psql_command" {
  value       = elestio_postgresql.postgres.database_admin.command
  description = "The PSQL command to connect to the database."
  sensitive   = true
}
```

- Run these commands in your terminal:

```bash
terraform init
terraform plan
terraform apply
eval "$(terraform output -raw psql_command)"
```

You have just deployed in a few lines of code a whole infrastructure.

## Sign-in limits

The Elestio API allows 15 sign-ins per hour per account. The provider signs in once per Terraform command and reuses the session token during that command. If you run many commands per hour on your own machine, set `ELESTIO_JWT_CACHE=on` to reuse the token between runs (it is stored in a private file in your user cache directory, so do not use it on shared CI runners). In CI, pass a token you obtained earlier with `ELESTIO_JWT`. See the [provider documentation](https://registry.terraform.io/providers/elestio/elestio/latest/docs#sign-in-limits-and-session-reuse) for details.

## License

terraform-provider-elestio is licensed under the MPL license. Full license text is available in the [LICENSE](LICENSE) file.
