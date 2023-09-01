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
      version = "0.2.0" # check out the latest version in the release section
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
  server_type   = "SMALL-1C-2G"
  provider_name = "hetzner"
  datacenter    = "fsn1"
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

## License

terraform-provider-elestio is licensed under the MPL license. Full license text is available in the [LICENSE](LICENSE) file.
