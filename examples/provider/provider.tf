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
  name             = "Demo"
  technical_emails = "admin@email.com"
}
