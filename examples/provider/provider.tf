terraform {
  required_providers {
    elestio = {
      source = "elestio/elestio"
    }
  }
}

# Configure the provider
provider "elestio" {
  email     = "<elestio_email>"
  api_token = "<elestio_api_token>"
}

# Create a project
resource "elestio_project" "project" {
  name = "Demo"
}

# Create a postgresql service
resource "elestio_postgresql" "database" {
  project_id    = elestio_project.project.id
  provider_name = "hetzner"
  datacenter    = "fsn1"
  server_type   = "MEDIUM-2C-4G"
}
