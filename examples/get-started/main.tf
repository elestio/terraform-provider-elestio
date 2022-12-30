terraform {
  required_providers {
    elestio = {
      source  = "elestio/elestio"
      version = "0.2.0" # check out the latest version in the release section
    }
  }
}

# Configure the Elestio Provider
provider "elestio" {
  email     = var.elestio_email
  api_token = var.elestio_api_token
}

# Create a Project
resource "elestio_project" "pg_project" {
  name             = "PostgreSQL Project"
  description      = "Contains a postgres database in Europe and Asia"
  technical_emails = var.elestio_email
}

# Create a PostgreSQL Service in Europe
resource "elestio_postgresql" "pg_europe" {
  project_id    = elestio_project.pg_project.id
  server_name   = "pg-europe"
  server_type   = "MICRO-1C-1G"
  provider_name = "lightsail"
  datacenter    = "eu-central-1"
  support_level = "level1"
  admin_email   = var.elestio_email
}

# Create a PostgreSQL Service in Asia
resource "elestio_postgresql" "pg_asia" {
  project_id    = elestio_project.pg_project.id
  server_name   = "pg-asia"
  server_type   = "MICRO-1C-1G"
  provider_name = "lightsail"
  datacenter    = "ap-northeast-2"
  support_level = "level1"
  admin_email   = var.elestio_email
}

# Extract the bash command to connect to the database with psql
output "pg_europe_psql" {
  value       = elestio_postgresql.pg_europe.database_admin.command
  description = "This is the bash command to connect to the europe database with psql"
  sensitive   = true
}

# Extract the bash command to connect to the database with psql
output "pg_asia_psql" {
  value       = elestio_postgresql.pg_asia.database_admin.command
  description = "This is the bash command to connect to the asia database with psql"
  sensitive   = true
}
