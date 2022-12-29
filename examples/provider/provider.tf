# Terraform 0.13+ requires providers to be declared in a "required_providers" block
terraform {
  required_providers {
    elestio = {
      source  = "elestio/elestio"
      version = "1.0.0"
    }
  }
}

# Configure the Elestio Provider
provider "elestio" {
  email     = "YOUR-EMAIL"
  api_token = "YOUR-API-TOKEN"
}

# Create a Project
resource "elestio_project" "myawesomeproject" {
  name = "Awesome project"

  # ...
}

